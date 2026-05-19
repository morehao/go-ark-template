package chunker

import (
	"strings"
	"unicode/utf8"
)

func init() {
	splitByHeadings = splitByHeadingsImpl
}

type headingBoundary struct {
	runeStart int
	line      string
}

func splitByHeadingsImpl(text string, cfg SplitterConfig, profile *DocProfile) []Chunk {
	if text == "" {
		return nil
	}
	if profile == nil {
		profile = ProfileDocument(text)
	}
	primaryLevel := profile.DominantHeadingLevel()
	if primaryLevel == 0 {
		return SplitText(text, cfg)
	}

	bounds := findHeadingBoundaries(text, primaryLevel)
	if len(bounds) <= 1 {
		return SplitText(text, cfg)
	}

	runes := []rune(text)
	hierarchy := NewHeadingHierarchy()

	var out []Chunk
	seq := 0

	for i, b := range bounds {
		endRune := len(runes)
		if i+1 < len(bounds) {
			endRune = bounds[i+1].runeStart
		}
		if b.line != "" {
			hierarchy.Observe(b.line)
		}
		breadcrumb := hierarchy.BreadcrumbWithHashes()
		observeSubHeadings(runes[b.runeStart:endRune], primaryLevel, hierarchy)

		sectionRunes := runes[b.runeStart:endRune]
		sectionContent := string(sectionRunes)
		secLen := len(sectionRunes)
		if secLen == 0 {
			continue
		}

		bcLen := utf8.RuneCountInString(breadcrumb)
		if bcLen+2+secLen <= cfg.ChunkSize {
			out = append(out, Chunk{
				Content:       sectionContent,
				ContextHeader: breadcrumb,
				Seq:           seq,
				Start:         b.runeStart,
				End:           endRune,
			})
			seq++
			continue
		}

		subChunks := SplitText(sectionContent, cfg)
		for _, sub := range subChunks {
			out = append(out, Chunk{
				Content:       sub.Content,
				ContextHeader: breadcrumb,
				Seq:           seq,
				Start:         b.runeStart + sub.Start,
				End:           b.runeStart + sub.End,
			})
			seq++
		}
	}

	return coalesceTinyChunks(out, cfg.ChunkSize)
}

func coalesceTinyChunks(in []Chunk, chunkSize int) []Chunk {
	if len(in) <= 1 || chunkSize <= 0 {
		return in
	}
	target := chunkSize / 2
	if target < 200 {
		target = 200
	}

	out := make([]Chunk, 0, len(in))
	cur := in[0]
	curLen := utf8.RuneCountInString(cur.Content)

	for i := 1; i < len(in); i++ {
		next := in[i]
		nextLen := utf8.RuneCountInString(next.Content)
		if cur.End == next.Start && curLen < target && curLen+nextLen <= chunkSize {
			cur.Content += next.Content
			cur.ContextHeader = commonHeadingPrefix(cur.ContextHeader, next.ContextHeader)
			cur.End = next.End
			curLen += nextLen
			continue
		}
		out = append(out, cur)
		cur = next
		curLen = nextLen
	}
	out = append(out, cur)

	for i := range out {
		out[i].Seq = i
	}
	return out
}

func commonHeadingPrefix(a, b string) string {
	if a == b {
		return a
	}
	la := strings.Split(a, "\n")
	lb := strings.Split(b, "\n")
	n := len(la)
	if len(lb) < n {
		n = len(lb)
	}
	common := 0
	for i := 0; i < n; i++ {
		if la[i] != lb[i] {
			break
		}
		common = i + 1
	}
	if common == 0 {
		return ""
	}
	return strings.Join(la[:common], "\n")
}

func findHeadingBoundaries(text string, primaryLevel int) []headingBoundary {
	runes := []rune(text)
	bounds := []headingBoundary{{runeStart: 0}}
	if len(runes) == 0 {
		return bounds
	}

	pos := 0
	inFence := false
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			pos += utf8.RuneCountInString(line)
			if i < len(lines)-1 {
				pos++
			}
			continue
		}
		if !inFence {
			m := MarkdownHeadingPattern.FindStringSubmatch(line)
			if m != nil {
				level := len(m[1])
				if level >= 1 && level <= primaryLevel && pos > 0 {
					bounds = append(bounds, headingBoundary{
						runeStart: pos,
						line:      line,
					})
				}
				if level >= 1 && level <= primaryLevel && pos == 0 {
					bounds[0].line = line
				}
			}
		}
		pos += utf8.RuneCountInString(line)
		if i < len(lines)-1 {
			pos++
		}
	}
	return bounds
}

func observeSubHeadings(runes []rune, primaryLevel int, h *HeadingHierarchy) {
	if len(runes) == 0 {
		return
	}
	text := string(runes)
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		m := MarkdownHeadingPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		level := len(m[1])
		if level > primaryLevel {
			h.Observe(line)
		}
	}
}
