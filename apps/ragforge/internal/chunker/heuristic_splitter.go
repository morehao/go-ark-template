package chunker

import (
	"sort"
	"strings"
	"unicode/utf8"
)

func init() {
	splitByHeuristics = splitByHeuristicsImpl
}

type boundary struct {
	runeStart int
	priority  int
}

func splitByHeuristicsImpl(text string, cfg SplitterConfig, _ *DocProfile) []Chunk {
	if text == "" {
		return nil
	}
	runes := []rune(text)
	totalRunes := len(runes)
	if totalRunes <= cfg.ChunkSize {
		return SplitText(text, cfg)
	}

	bounds := findHeuristicBoundaries(text, cfg.Languages)
	if prot := protectedSpansRune(text, protectedSpans(text)); len(prot) > 0 {
		bounds = dropBoundsInsideSpans(bounds, prot)
	}
	if len(bounds) == 0 {
		return SplitText(text, cfg)
	}

	bounds = append(bounds, boundary{runeStart: totalRunes})
	if bounds[0].runeStart != 0 {
		bounds = append([]boundary{{runeStart: 0}}, bounds...)
	}

	var out []Chunk
	seq := 0
	chunkStart := bounds[0].runeStart
	curEnd := chunkStart
	minChunkSize := cfg.ChunkSize / 4
	if minChunkSize < 50 {
		minChunkSize = 50
	}

	for i := 1; i < len(bounds); i++ {
		nextEnd := bounds[i].runeStart
		blockLen := nextEnd - curEnd

		if blockLen > cfg.ChunkSize {
			if curEnd-chunkStart > 0 {
				out = appendChunk(out, runes, chunkStart, curEnd, &seq)
			}
			out = appendOversizeBlock(out, runes, curEnd, nextEnd, cfg, &seq)
			curEnd = nextEnd
			chunkStart = nextEnd
			continue
		}

		accumulated := nextEnd - chunkStart
		if accumulated > cfg.ChunkSize && curEnd-chunkStart >= minChunkSize {
			out = appendChunk(out, runes, chunkStart, curEnd, &seq)
			chunkStart = applyOverlapAligned(runes, curEnd, cfg.ChunkOverlap, bounds)
		}
		curEnd = nextEnd
	}

	if curEnd > chunkStart {
		out = appendChunk(out, runes, chunkStart, curEnd, &seq)
	}
	return out
}

func findHeuristicBoundaries(text string, langs []string) []boundary {
	var bounds []boundary

	for _, idx := range allRuneIndices(text, "\f") {
		bounds = append(bounds, boundary{runeStart: idx, priority: PrioFormFeed})
	}

	lines := strings.Split(text, "\n")
	chapterPatterns := ChapterPatternsForLangs(langs)
	pos := 0
	inFence := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
		} else if !inFence {
			runeStart := pos
			added := false
			for _, pat := range chapterPatterns {
				if pat.MatchString(line) {
					bounds = append(bounds, boundary{runeStart: runeStart, priority: PrioChapterMarker})
					added = true
					break
				}
			}
			if !added && NumberedSectionPattern.MatchString(line) {
				bounds = append(bounds, boundary{runeStart: runeStart, priority: PrioNumberedHead})
				added = true
			}
			if !added && AllCapsHeadingPattern.MatchString(line) {
				bounds = append(bounds, boundary{runeStart: runeStart, priority: PrioAllCapsHeading})
				added = true
			}
			if !added && VisualSeparatorPattern.MatchString(line) {
				bounds = append(bounds, boundary{runeStart: runeStart, priority: PrioVisualSep})
				added = true
			}
			if !added && PageFooterPattern.MatchString(line) {
				bounds = append(bounds, boundary{runeStart: runeStart, priority: PrioPageFooter})
			}
		}
		pos += utf8.RuneCountInString(line)
		if i < len(lines)-1 {
			pos++
		}
	}

	for _, idx := range ExcessiveBlanksPattern.FindAllStringIndex(text, -1) {
		runeStart := utf8.RuneCountInString(text[:idx[1]])
		bounds = append(bounds, boundary{runeStart: runeStart, priority: PrioBlankBlock})
	}

	if len(bounds) == 0 {
		return nil
	}

	sort.Slice(bounds, func(i, j int) bool {
		if bounds[i].runeStart != bounds[j].runeStart {
			return bounds[i].runeStart < bounds[j].runeStart
		}
		return bounds[i].priority > bounds[j].priority
	})
	deduped := bounds[:0]
	prev := -1
	for _, b := range bounds {
		if b.runeStart != prev {
			deduped = append(deduped, b)
			prev = b.runeStart
		}
	}
	return deduped
}

func dropBoundsInsideSpans(bounds []boundary, spans []span) []boundary {
	if len(spans) == 0 {
		return bounds
	}
	out := bounds[:0]
boundLoop:
	for _, b := range bounds {
		for _, s := range spans {
			if s.start >= b.runeStart {
				break
			}
			if b.runeStart < s.end {
				continue boundLoop
			}
		}
		out = append(out, b)
	}
	return out
}

func allRuneIndices(text, needle string) []int {
	var out []int
	if needle == "" {
		return out
	}
	pos := 0
	for _, r := range text {
		if string(r) == needle {
			out = append(out, pos)
		}
		pos++
	}
	return out
}

func appendChunk(out []Chunk, runes []rune, start, end int, seq *int) []Chunk {
	if end <= start {
		return out
	}
	raw := string(runes[start:end])
	if strings.TrimSpace(raw) == "" {
		return out
	}
	c := Chunk{Content: raw, Seq: *seq, Start: start, End: end}
	*seq++
	return append(out, c)
}

func appendOversizeBlock(out []Chunk, runes []rune, start, end int, cfg SplitterConfig, seq *int) []Chunk {
	if end <= start {
		return out
	}
	subText := string(runes[start:end])
	subs := SplitText(subText, cfg)
	for _, s := range subs {
		out = append(out, Chunk{
			Content: s.Content,
			Seq:     *seq,
			Start:   start + s.Start,
			End:     start + s.End,
		})
		*seq++
	}
	return out
}

func applyOverlapAligned(runes []rune, curEnd, overlap int, bounds []boundary) int {
	if overlap <= 0 {
		return curEnd
	}
	target := curEnd - overlap
	if target < 0 {
		target = 0
	}
	windowStart := curEnd - 2*overlap
	if windowStart < 0 {
		windowStart = 0
	}

	bestBound := -1
	for _, b := range bounds {
		if b.runeStart >= windowStart && b.runeStart < curEnd && b.runeStart > bestBound {
			bestBound = b.runeStart
		}
	}
	if bestBound >= 0 {
		return bestBound
	}

	for i := target; i > windowStart && i < len(runes); i-- {
		if runes[i] == '\n' {
			return i + 1
		}
	}
	return target
}
