package chunker

import (
	"regexp"
	"sort"
	"strings"
)

type headerTrackerHook struct {
	startPattern *regexp.Regexp
	endPattern   *regexp.Regexp
	priority     int
}

var defaultHeaderHooks = []headerTrackerHook{
	{
		startPattern: regexp.MustCompile(`(?si)^\s*(?:\|[^|\n]*)+[\r\n]+\s*(?:\|\s*:?-{3,}:?\s*)+\|?[\r\n]+$`),
		endPattern:   regexp.MustCompile(`(?si)^\s*$|^\s*[^|\s].*$`),
		priority:     15,
	},
}

var tableRowPattern = regexp.MustCompile(`(?m)^\s*(?:\|[^|\n]*)+\|\s*$`)

type headerTracker struct {
	hooks         []headerTrackerHook
	activeHeaders map[int]string
	endedHeaders  map[int]bool
	pendingExtend map[int]bool
}

func newHeaderTracker() *headerTracker {
	return &headerTracker{
		hooks:         defaultHeaderHooks,
		activeHeaders: make(map[int]string),
		endedHeaders:  make(map[int]bool),
		pendingExtend: make(map[int]bool),
	}
}

func (ht *headerTracker) update(split string) {
	for _, hook := range ht.hooks {
		if _, active := ht.activeHeaders[hook.priority]; active {
			if hook.endPattern.MatchString(split) {
				ht.endedHeaders[hook.priority] = true
				delete(ht.activeHeaders, hook.priority)
				delete(ht.pendingExtend, hook.priority)
			}
		}
	}

	for p := range ht.pendingExtend {
		if _, active := ht.activeHeaders[p]; active && tableRowPattern.MatchString(split) {
			sep := extractSeparatorLine(ht.activeHeaders[p])
			ht.activeHeaders[p] = split + sep
		}
		delete(ht.pendingExtend, p)
	}

	for _, hook := range ht.hooks {
		if _, active := ht.activeHeaders[hook.priority]; active {
			continue
		}
		if ht.endedHeaders[hook.priority] {
			continue
		}
		if loc := hook.startPattern.FindString(split); loc != "" {
			ht.activeHeaders[hook.priority] = loc
			if isEmptyTableHeaderRow(loc) {
				ht.pendingExtend[hook.priority] = true
			}
		}
	}

	if len(ht.activeHeaders) == 0 {
		for k := range ht.endedHeaders {
			delete(ht.endedHeaders, k)
		}
	}
}

func (ht *headerTracker) getHeaders() string {
	if len(ht.activeHeaders) == 0 {
		return ""
	}

	type entry struct {
		priority int
		text     string
	}
	entries := make([]entry, 0, len(ht.activeHeaders))
	for p, t := range ht.activeHeaders {
		entries = append(entries, entry{p, t})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].priority > entries[j].priority
	})

	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = e.text
	}
	return strings.Join(parts, "\n")
}

func isEmptyTableHeaderRow(header string) bool {
	idx := strings.IndexByte(header, '\n')
	if idx < 0 {
		return false
	}
	row := strings.TrimSpace(header[:idx])
	for _, r := range row {
		if r != '|' && r != ' ' && r != '\t' {
			return false
		}
	}
	return true
}

func extractSeparatorLine(header string) string {
	for _, line := range strings.Split(header, "\n") {
		if strings.Contains(line, "---") {
			return line + "\n"
		}
	}
	return ""
}
