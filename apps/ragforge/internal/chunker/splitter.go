package chunker

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var protectedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?s)\$\$.*?\$\$`),
	regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`),
	regexp.MustCompile(`\[[^\]]*\]\([^)]+\)`),
	regexp.MustCompile(`(?m)[ ]*(?:\|[^|\n]*)+\|[\r\n]+\s*(?:\|\s*:?-{3,}:?\s*)+\|[\r\n]+`),
	regexp.MustCompile(`(?m)[ ]*(?:\|[^|\n]*)+\|[\r\n]+`),
	regexp.MustCompile("(?s)```(?:\\w+)?[\\r\\n].*?```"),
}

type span struct {
	start, end int
}

func protectedSpansRune(text string, byteSpans []span) []span {
	if len(byteSpans) == 0 {
		return nil
	}
	out := make([]span, 0, len(byteSpans))
	runeIdx := 0
	byteIdx := 0
	for _, s := range byteSpans {
		for byteIdx < s.start && byteIdx < len(text) {
			_, size := utf8.DecodeRuneInString(text[byteIdx:])
			byteIdx += size
			runeIdx++
		}
		startRune := runeIdx
		for byteIdx < s.end && byteIdx < len(text) {
			_, size := utf8.DecodeRuneInString(text[byteIdx:])
			byteIdx += size
			runeIdx++
		}
		out = append(out, span{start: startRune, end: runeIdx})
	}
	return out
}

func protectedSpans(text string) []span {
	type match struct {
		start, end int
	}
	var all []match
	for _, pat := range protectedPatterns {
		locs := pat.FindAllStringIndex(text, -1)
		for _, loc := range locs {
			if loc[1]-loc[0] > 0 {
				all = append(all, match{loc[0], loc[1]})
			}
		}
	}
	if len(all) == 0 {
		return nil
	}

	for i := 1; i < len(all); i++ {
		for j := i; j > 0; j-- {
			if all[j].start < all[j-1].start ||
				(all[j].start == all[j-1].start && (all[j].end-all[j].start) > (all[j-1].end-all[j-1].start)) {
				all[j], all[j-1] = all[j-1], all[j]
			} else {
				break
			}
		}
	}

	var result []span
	lastEnd := 0
	for _, m := range all {
		if m.start >= lastEnd {
			result = append(result, span(m))
			lastEnd = m.end
		}
	}
	return result
}

type splitUnit struct {
	text       string
	start, end int
}

func splitBySeparators(text string, separators []string, chunkSize int) []string {
	if text == "" || len(separators) == 0 {
		return []string{text}
	}
	if chunkSize > 0 && runeLen(text) <= chunkSize {
		return []string{text}
	}

	for i, sep := range separators {
		if sep == "" {
			continue
		}
		re := regexp.MustCompile("(" + regexp.QuoteMeta(sep) + ")")
		splits := re.Split(text, -1)
		matches := re.FindAllString(text, -1)
		if len(matches) == 0 {
			continue
		}

		var pieces []string
		for j, s := range splits {
			if s != "" {
				pieces = append(pieces, s)
			}
			if j < len(matches) && matches[j] != "" {
				pieces = append(pieces, matches[j])
			}
		}
		if len(pieces) <= 1 {
			continue
		}

		var out []string
		remaining := separators[i+1:]
		for _, p := range pieces {
			if chunkSize > 0 && runeLen(p) > chunkSize && len(remaining) > 0 {
				out = append(out, splitBySeparators(p, remaining, chunkSize)...)
			} else {
				out = append(out, p)
			}
		}
		return out
	}
	return []string{text}
}

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

// SplitText splits text into chunks with overlap, respecting protected patterns.
func SplitText(text string, cfg SplitterConfig) []Chunk {
	if text == "" {
		return nil
	}

	chunkSize := cfg.ChunkSize
	chunkOverlap := cfg.ChunkOverlap
	separators := cfg.Separators

	if chunkSize <= 0 {
		chunkSize = 512
	}
	if chunkOverlap < 0 {
		chunkOverlap = 0
	}

	protected := protectedSpans(text)

	units := buildUnitsWithProtection(text, protected, separators, chunkSize)

	return mergeUnits(units, chunkSize, chunkOverlap)
}

func buildUnitsWithProtection(text string, protected []span, separators []string, chunkSize int) []splitUnit {
	const maxProtectedSize = 7500

	var units []splitUnit
	bytePos := 0
	runePos := 0

	for _, p := range protected {
		if p.start > bytePos {
			pre := text[bytePos:p.start]
			parts := splitBySeparators(pre, separators, chunkSize)
			runeOffset := runePos
			for _, part := range parts {
				partRuneLen := runeLen(part)
				units = append(units, splitUnit{
					text:  part,
					start: runeOffset,
					end:   runeOffset + partRuneLen,
				})
				runeOffset += partRuneLen
			}
			runePos += runeLen(pre)
		}

		protText := text[p.start:p.end]
		protRuneLen := runeLen(protText)

		if protRuneLen > maxProtectedSize {
			runes := []rune(protText)
			offset := 0
			for offset < len(runes) {
				chunkEnd := offset + maxProtectedSize
				if chunkEnd > len(runes) {
					chunkEnd = len(runes)
				} else {
					for i := chunkEnd - 1; i > offset && i > chunkEnd-200; i-- {
						if runes[i] == '\n' || runes[i] == ' ' {
							chunkEnd = i + 1
							break
						}
					}
				}

				chunkText := string(runes[offset:chunkEnd])
				chunkLen := chunkEnd - offset
				units = append(units, splitUnit{
					text:  chunkText,
					start: runePos + offset,
					end:   runePos + offset + chunkLen,
				})
				offset = chunkEnd
			}
		} else {
			units = append(units, splitUnit{
				text:  protText,
				start: runePos,
				end:   runePos + protRuneLen,
			})
		}
		runePos += protRuneLen
		bytePos = p.end
	}

	if bytePos < len(text) {
		remaining := text[bytePos:]
		parts := splitBySeparators(remaining, separators, chunkSize)
		runeOffset := runePos
		for _, part := range parts {
			partRuneLen := runeLen(part)
			units = append(units, splitUnit{
				text:  part,
				start: runeOffset,
				end:   runeOffset + partRuneLen,
			})
			runeOffset += partRuneLen
		}
	}

	return units
}

func mergeUnits(units []splitUnit, chunkSize, chunkOverlap int) []Chunk {
	if len(units) == 0 {
		return nil
	}

	const absoluteMaxSize = 7500

	ht := newHeaderTracker()

	var chunks []Chunk
	var current []splitUnit
	curLen := 0

	for _, u := range units {
		uLen := runeLen(u.text)

		if uLen > absoluteMaxSize {
			if len(current) > 0 {
				chunks = append(chunks, buildChunk(current, len(chunks)))
				current = nil
				curLen = 0
			}

			ht.update(u.text)

			runes := []rune(u.text)
			offset := 0
			for offset < len(runes) {
				chunkEnd := offset + absoluteMaxSize
				if chunkEnd > len(runes) {
					chunkEnd = len(runes)
				} else {
					for i := chunkEnd - 1; i > offset && i > chunkEnd-200; i-- {
						if runes[i] == '\n' || runes[i] == ' ' {
							chunkEnd = i + 1
							break
						}
					}
				}

				chunkText := string(runes[offset:chunkEnd])
				chunks = append(chunks, Chunk{
					Content: chunkText,
					Seq:     len(chunks),
					Start:   u.start + offset,
					End:     u.start + chunkEnd,
				})
				offset = chunkEnd
			}
			continue
		}

		ht.update(u.text)
		headers := ht.getHeaders()
		headersLen := runeLen(headers)
		if headersLen > chunkSize {
			headers = ""
			headersLen = 0
		}

		if curLen+uLen+headersLen > chunkSize && len(current) > 0 {
			chunks = append(chunks, buildChunk(current, len(chunks)))

			current, curLen = computeOverlap(current, chunkOverlap, chunkSize, uLen)

			if headers != "" && headersLen+uLen <= chunkSize {
				for len(current) > 0 && curLen+uLen+headersLen > chunkSize {
					curLen -= runeLen(current[0].text)
					current = current[1:]
				}

				overlapText := unitsText(current)
				if !headerAlreadyPresent(headers, overlapText, u.text) {
					startPos := u.start
					if len(current) > 0 {
						startPos = current[0].start
					}
					hUnit := splitUnit{text: headers, start: startPos, end: startPos}
					current = append([]splitUnit{hUnit}, current...)
					curLen += headersLen
				}
			}
		}

		if curLen+uLen > absoluteMaxSize {
			if len(current) > 0 {
				chunks = append(chunks, buildChunk(current, len(chunks)))
				current = nil
				curLen = 0
			}
		}

		current = append(current, u)
		curLen += uLen
	}

	if len(current) > 0 {
		chunks = append(chunks, buildChunk(current, len(chunks)))
	}

	return chunks
}

func unitsText(units []splitUnit) string {
	var sb strings.Builder
	for _, u := range units {
		sb.WriteString(u.text)
	}
	return sb.String()
}

func headerAlreadyPresent(headers, overlapText, unitText string) bool {
	if strings.Contains(overlapText, headers) || strings.Contains(unitText, headers) {
		return true
	}

	colRow := headerColumnRow(headers)
	if colRow == "" {
		return false
	}

	return strings.Contains(overlapText, colRow) || strings.Contains(unitText, colRow)
}

func headerColumnRow(header string) string {
	for _, line := range strings.Split(header, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "---") {
			continue
		}
		onlyPipes := true
		for _, r := range line {
			if r != '|' && r != ' ' && r != '\t' {
				onlyPipes = false
				break
			}
		}
		if !onlyPipes {
			return line
		}
	}
	return ""
}

func buildChunk(units []splitUnit, seq int) Chunk {
	var sb strings.Builder
	for _, u := range units {
		sb.WriteString(u.text)
	}
	return Chunk{
		Content: sb.String(),
		Seq:     seq,
		Start:   units[0].start,
		End:     units[len(units)-1].end,
	}
}

func computeOverlap(current []splitUnit, chunkOverlap, chunkSize, nextLen int) ([]splitUnit, int) {
	if chunkOverlap <= 0 {
		return nil, 0
	}

	overlapLen := 0
	startIdx := len(current)
	for i := len(current) - 1; i >= 0; i-- {
		uLen := runeLen(current[i].text)
		if overlapLen+uLen > chunkOverlap {
			break
		}
		if overlapLen+uLen+nextLen > chunkSize {
			break
		}
		overlapLen += uLen
		startIdx = i
	}

	for startIdx < len(current) {
		u := current[startIdx]
		isHeaderMarker := u.start == u.end
		trimmed := strings.TrimSpace(u.text)
		if isHeaderMarker || trimmed == "" || isSeparatorOnly(u.text) {
			overlapLen -= runeLen(u.text)
			startIdx++
		} else {
			break
		}
	}

	if startIdx >= len(current) {
		return nil, 0
	}

	overlap := make([]splitUnit, len(current)-startIdx)
	copy(overlap, current[startIdx:])
	return overlap, overlapLen
}

func isSeparatorOnly(s string) bool {
	for _, r := range s {
		if r != '\n' && r != '\r' && r != ' ' && r != '\t' && r != '。' {
			return false
		}
	}
	return true
}

// SplitTextParentChild performs two-level chunking:
//  1. Split text into large parent chunks (parentCfg).
//  2. Split each parent into smaller child chunks (childCfg) for embedding.
func SplitTextParentChild(text string, parentCfg, childCfg SplitterConfig) ParentChildResult {
	parents := SplitText(text, parentCfg)
	if len(parents) == 0 {
		return ParentChildResult{}
	}

	var newParents []Chunk
	var children []ChildChunk
	childSeq := 0
	for _, parent := range parents {
		subs := SplitText(parent.Content, childCfg)

		parentIndex := -1
		if len(subs) > 1 || (len(subs) == 1 && subs[0].Content != parent.Content) {
			parentIndex = len(newParents)
			newParents = append(newParents, parent)
		}

		for _, sub := range subs {
			sub.Seq = childSeq
			sub.Start += parent.Start
			sub.End += parent.Start
			children = append(children, ChildChunk{
				Chunk:       sub,
				ParentIndex: parentIndex,
			})
			childSeq++
		}
	}
	return ParentChildResult{Parents: newParents, Children: children}
}

// ExtractImageRefs extracts markdown image references from text.
var imageRefPattern = regexp.MustCompile(`!\[([^\]]*)\]\(([^()\s]*(?:\([^)]*\)[^()\s]*)*)\)`)

func ExtractImageRefs(text string) []ImageRef {
	matches := imageRefPattern.FindAllStringSubmatchIndex(text, -1)
	var refs []ImageRef
	for _, m := range matches {
		refs = append(refs, ImageRef{
			OriginalRef: text[m[4]:m[5]],
			AltText:     text[m[2]:m[3]],
			Start:       m[0],
			End:         m[1],
		})
	}
	return refs
}
