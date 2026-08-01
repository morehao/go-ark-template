package chunker

import "strings"

// HeadingHierarchy maintains a stack of active Markdown headings indexed by
// level (1..6).
type HeadingHierarchy struct {
	stack [6]string
	depth int
}

// NewHeadingHierarchy returns an empty hierarchy.
func NewHeadingHierarchy() *HeadingHierarchy {
	return &HeadingHierarchy{}
}

// Observe parses line and updates the hierarchy if line is a Markdown heading.
func (h *HeadingHierarchy) Observe(line string) (int, string) {
	m := MarkdownHeadingPattern.FindStringSubmatch(line)
	if m == nil {
		return 0, ""
	}
	level := len(m[1])
	if level < 1 || level > 6 {
		return 0, ""
	}
	heading := strings.TrimSpace(m[2])
	h.stack[level-1] = heading
	for i := level; i < 6; i++ {
		h.stack[i] = ""
	}
	if level > h.depth {
		h.depth = level
	} else {
		h.depth = 0
		for i := 0; i < 6; i++ {
			if h.stack[i] != "" {
				h.depth = i + 1
			}
		}
	}
	return level, heading
}

// Breadcrumb returns the current heading path joined by " > ".
func (h *HeadingHierarchy) Breadcrumb() string {
	if h.depth == 0 {
		return ""
	}
	parts := make([]string, 0, h.depth)
	for i := 0; i < h.depth; i++ {
		if h.stack[i] != "" {
			parts = append(parts, h.stack[i])
		}
	}
	return strings.Join(parts, " > ")
}

// BreadcrumbWithHashes returns the path with the original `#` prefixes.
func (h *HeadingHierarchy) BreadcrumbWithHashes() string {
	if h.depth == 0 {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < h.depth; i++ {
		if h.stack[i] == "" {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strings.Repeat("#", i+1))
		sb.WriteByte(' ')
		sb.WriteString(h.stack[i])
	}
	return sb.String()
}

// Depth returns the current deepest active heading level.
func (h *HeadingHierarchy) Depth() int { return h.depth }

// Reset clears all state.
func (h *HeadingHierarchy) Reset() {
	for i := range h.stack {
		h.stack[i] = ""
	}
	h.depth = 0
}
