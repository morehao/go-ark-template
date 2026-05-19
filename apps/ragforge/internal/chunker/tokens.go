package chunker

import (
	"unicode"
	"unicode/utf8"
)

// Language identifiers used by the token estimator and the heuristic splitter.
const (
	LangEnglish = "en"
	LangGerman  = "de"
	LangChinese = "zh"
	LangMixed   = "mixed"
)

var charsPerToken = map[string]float64{
	LangEnglish: 4.0,
	LangGerman:  4.5,
	LangChinese: 1.7,
	LangMixed:   3.0,
}

// ApproxTokenCount returns a conservative token estimate for s in the given language.
func ApproxTokenCount(s string, lang string) int {
	if s == "" {
		return 0
	}
	return ApproxTokenCountFromRuneLen(utf8.RuneCountInString(s), lang)
}

// ApproxTokenCountFromRuneLen is the allocation-free variant.
func ApproxTokenCountFromRuneLen(runeLen int, lang string) int {
	if runeLen <= 0 {
		return 0
	}
	ratio, ok := charsPerToken[lang]
	if !ok {
		ratio = charsPerToken[LangMixed]
	}
	approx := float64(runeLen) / ratio
	if approx < 1 {
		return 1
	}
	return int(approx + 0.5)
}

// DetectLanguage returns a coarse language label by counting CJK runes vs. Latin runes.
func DetectLanguage(s string) string {
	if s == "" {
		return LangMixed
	}
	var cjk, latin, umlaut int
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hangul, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r):
			cjk++
		case isGermanUmlaut(r):
			umlaut++
			latin++
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			latin++
		}
	}
	total := cjk + latin
	if total == 0 {
		return LangMixed
	}
	cjkRatio := float64(cjk) / float64(total)
	latinRatio := float64(latin) / float64(total)
	if cjkRatio >= 0.15 && latinRatio >= 0.15 {
		return LangMixed
	}
	if cjkRatio > 0.3 {
		return LangChinese
	}
	if umlaut > 0 || hasGermanWords(s) {
		return LangGerman
	}
	return LangEnglish
}

func isGermanUmlaut(r rune) bool {
	switch r {
	case 'ä', 'ö', 'ü', 'Ä', 'Ö', 'Ü', 'ß':
		return true
	}
	return false
}

func hasGermanWords(s string) bool {
	const sample = 512
	if len(s) > sample {
		s = s[:sample]
	}
	for _, w := range []string{" der ", " die ", " das ", " und ", " ist ", " nicht ", " mit ", " auf "} {
		if containsLower(s, w) {
			return true
		}
	}
	return false
}

func containsLower(haystack, needle string) bool {
	if len(haystack) < len(needle) {
		return false
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		match := true
		for j := 0; j < len(needle); j++ {
			h := haystack[i+j]
			if h >= 'A' && h <= 'Z' {
				h += 'a' - 'A'
			}
			if h != needle[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// CharsForTokenLimit converts a token limit into an approximate character budget.
func CharsForTokenLimit(tokens int, lang string) int {
	if tokens <= 0 {
		return 0
	}
	ratio, ok := charsPerToken[lang]
	if !ok {
		ratio = charsPerToken[LangMixed]
	}
	return int(float64(tokens) * ratio * 0.9)
}
