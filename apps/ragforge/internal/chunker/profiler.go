package chunker

import (
	"math"
	"strings"
)

// DocProfile holds the document-level signals used to choose a chunking tier.
type DocProfile struct {
	TotalChars int     `json:"total_chars"`
	TotalLines int     `json:"total_lines"`
	AvgLineLen float64 `json:"avg_line_len"`
	StdLineLen float64 `json:"std_line_len"`

	MdHeadingCounts map[int]int `json:"md_heading_counts"`
	MdHeadingTotal  int         `json:"md_heading_total"`

	NumberedSectionCount  int `json:"numbered_section_count"`
	AllCapsShortLineCount int `json:"all_caps_short_line_count"`
	BlankParagraphBreaks  int `json:"blank_paragraph_breaks"`
	FormFeedCount         int `json:"form_feed_count"`
	VisualSepCount        int `json:"visual_sep_count"`
	GermanChapterCount    int `json:"german_chapter_count"`
	EnglishChapterCount   int `json:"english_chapter_count"`
	ChineseChapterCount   int `json:"chinese_chapter_count"`
	RepeatedFooterCount   int `json:"repeated_footer_count"`

	HasTables bool    `json:"has_tables"`
	HasCode   bool    `json:"has_code"`
	CodeRatio float64 `json:"code_ratio"`

	DetectedLangs []string `json:"detected_langs"`
}

// HeadingDensity returns the share of lines that are Markdown headings.
func (p *DocProfile) HeadingDensity() float64 {
	if p.TotalLines == 0 {
		return 0
	}
	return float64(p.MdHeadingTotal) / float64(p.TotalLines)
}

// DominantHeadingLevel returns the heading level (1..6) that should drive section splitting.
func (p *DocProfile) DominantHeadingLevel() int {
	if p.MdHeadingTotal == 0 {
		return 0
	}
	for level := 1; level <= 6; level++ {
		if p.MdHeadingCounts[level] >= 3 {
			return level
		}
	}
	for level := 6; level >= 1; level-- {
		if p.MdHeadingCounts[level] > 0 {
			return level
		}
	}
	return 0
}

// HeuristicMarkerTotal sums the non-Markdown structural markers.
func (p *DocProfile) HeuristicMarkerTotal() int {
	return p.NumberedSectionCount +
		p.GermanChapterCount + p.EnglishChapterCount + p.ChineseChapterCount +
		p.AllCapsShortLineCount + p.VisualSepCount + p.FormFeedCount
}

// ProfileDocument runs a single pass over text and returns its profile.
func ProfileDocument(text string) *DocProfile {
	p := &DocProfile{
		MdHeadingCounts: make(map[int]int),
	}
	if text == "" {
		return p
	}

	p.TotalChars = len([]rune(text))
	p.FormFeedCount = strings.Count(text, "\f")

	lines := strings.Split(text, "\n")
	p.TotalLines = len(lines)

	var lengths []float64
	inFence := false
	codeChars := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			p.HasCode = true
			continue
		}
		if inFence {
			codeChars += len([]rune(line))
			continue
		}

		runeLen := len([]rune(line))
		lengths = append(lengths, float64(runeLen))

		if matchHeading(line, &p.MdHeadingCounts) {
			p.MdHeadingTotal++
			continue
		}
		if NumberedSectionPattern.MatchString(line) {
			p.NumberedSectionCount++
		}
		if GermanChapterPattern.MatchString(line) {
			p.GermanChapterCount++
		}
		if EnglishChapterPattern.MatchString(line) {
			p.EnglishChapterCount++
		}
		if ChineseChapterPattern.MatchString(line) {
			p.ChineseChapterCount++
		}
		if AllCapsHeadingPattern.MatchString(line) {
			p.AllCapsShortLineCount++
		}
		if VisualSeparatorPattern.MatchString(line) {
			p.VisualSepCount++
		}
		if PageFooterPattern.MatchString(line) {
			p.RepeatedFooterCount++
		}
		if strings.HasPrefix(trimmed, "|") && strings.HasSuffix(trimmed, "|") {
			p.HasTables = true
		}
	}

	if len(lengths) > 0 {
		var sum float64
		for _, l := range lengths {
			sum += l
		}
		p.AvgLineLen = sum / float64(len(lengths))
		var variance float64
		for _, l := range lengths {
			d := l - p.AvgLineLen
			variance += d * d
		}
		variance /= float64(len(lengths))
		p.StdLineLen = math.Sqrt(variance)
	}

	if p.TotalChars > 0 {
		p.CodeRatio = float64(codeChars) / float64(p.TotalChars)
	}

	p.BlankParagraphBreaks = strings.Count(text, "\n\n\n")

	sample := text
	if len(sample) > 4096 {
		sample = sample[:4096]
	}
	lang := DetectLanguage(sample)
	p.DetectedLangs = []string{lang}
	if lang == LangMixed {
		p.DetectedLangs = []string{LangEnglish, LangGerman, LangChinese}
	}

	return p
}

func matchHeading(line string, counts *map[int]int) bool {
	m := MarkdownHeadingPattern.FindStringSubmatch(line)
	if m == nil {
		return false
	}
	level := len(m[1])
	if level < 1 || level > 6 {
		return false
	}
	(*counts)[level]++
	return true
}

// StrategyTier identifies which chunking implementation should run.
type StrategyTier string

const (
	TierHeading   StrategyTier = "heading"
	TierHeuristic StrategyTier = "heuristic"
	TierLegacy    StrategyTier = "legacy"
)

// SelectStrategy returns the ordered tier chain to attempt for this document.
func SelectStrategy(p *DocProfile) []StrategyTier {
	if p == nil {
		return []StrategyTier{TierLegacy}
	}
	var chain []StrategyTier

	if p.MdHeadingTotal >= 3 && p.HeadingDensity() > 0.005 && p.DominantHeadingLevel() > 0 {
		chain = append(chain, TierHeading)
	}

	if p.HeuristicMarkerTotal() >= 5 || p.FormFeedCount > 0 ||
		p.GermanChapterCount+p.EnglishChapterCount+p.ChineseChapterCount > 0 {
		chain = append(chain, TierHeuristic)
	}

	chain = append(chain, TierLegacy)
	return chain
}
