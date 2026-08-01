package chunker

import "math"

// ValidationResult captures the verdict and reason for a chunk-set.
type ValidationResult struct {
	OK     bool
	Reason string
}

// ValidateChunks checks whether the given chunks form a usable result.
func ValidateChunks(chunks []Chunk, totalChars, chunkSize int) ValidationResult {
	if len(chunks) == 0 {
		return ValidationResult{Reason: "no chunks produced"}
	}

	if len(chunks) == 1 && totalChars > 2*chunkSize {
		return ValidationResult{Reason: "single chunk for large document"}
	}

	var sum, sumSq float64
	maxLen, minLen := 0, math.MaxInt32
	for _, c := range chunks {
		l := len([]rune(c.Content))
		sum += float64(l)
		sumSq += float64(l * l)
		if l > maxLen {
			maxLen = l
		}
		if l < minLen {
			minLen = l
		}
	}
	avg := sum / float64(len(chunks))

	tinyCount := 0
	for i, c := range chunks {
		if i == len(chunks)-1 {
			continue
		}
		if len([]rune(c.Content)) < 50 {
			tinyCount++
		}
	}
	if tinyCount > len(chunks)/4 && tinyCount > 2 {
		return ValidationResult{Reason: "too many tiny chunks"}
	}

	if maxLen < chunkSize/4 && totalChars > chunkSize {
		return ValidationResult{Reason: "all chunks far below target size"}
	}

	if maxLen > 2*chunkSize && chunkSize > 0 {
		return ValidationResult{Reason: "chunk exceeds 2x target size"}
	}

	_ = avg
	return ValidationResult{OK: true}
}
