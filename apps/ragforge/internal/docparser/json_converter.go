package docparser

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const defaultJSONChunkSize = 1536

var minJSONChunkSize = defaultJSONChunkSize - 200

func jsonToMarkdown(data []byte) (string, error) {
	data = trimBOM(data)
	if len(data) == 0 {
		return "", fmt.Errorf("empty JSON content")
	}
	if !json.Valid(data) {
		return "", fmt.Errorf("invalid JSON content")
	}

	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	normalized := listToDictPreprocess(parsed)

	wholeSize := jsonSize(normalized)
	if wholeSize <= defaultJSONChunkSize {
		formatted := formatValue(normalized)
		return wrapCodeBlock(formatted), nil
	}

	chunks := recursiveJSONSplit(normalized, nil, nil)

	blocks := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}
		blocks = append(blocks, wrapCodeBlock(formatValue(chunk)))
	}

	if len(blocks) == 0 {
		return wrapCodeBlock(formatValue(normalized)), nil
	}
	return strings.Join(blocks, "\n\n"), nil
}

func recursiveJSONSplit(
	data interface{},
	currentPath []string,
	chunks []map[string]interface{},
) []map[string]interface{} {
	if chunks == nil {
		chunks = []map[string]interface{}{{}}
	}

	dict, ok := data.(map[string]interface{})
	if !ok {
		if len(currentPath) > 0 && len(chunks) > 0 {
			setNestedDict(chunks[len(chunks)-1], currentPath, data)
		}
		return chunks
	}

	keys := sortedKeys(dict)

	for _, key := range keys {
		value := dict[key]
		newPath := append(append([]string{}, currentPath...), key)

		chunkSize := jsonSize(chunks[len(chunks)-1])
		itemSize := jsonSize(map[string]interface{}{key: value})
		remaining := defaultJSONChunkSize - chunkSize

		if itemSize <= remaining {
			setNestedDict(chunks[len(chunks)-1], newPath, value)
		} else {
			if chunkSize >= minJSONChunkSize {
				chunks = append(chunks, map[string]interface{}{})
			}

			normalized := listToDictPreprocess(value)
			if subDict, isDict := normalized.(map[string]interface{}); isDict && canSplitDict(subDict) {
				chunks = recursiveJSONSplit(subDict, newPath, chunks)
			} else {
				setNestedDict(chunks[len(chunks)-1], newPath, value)
			}
		}
	}

	return chunks
}

func setNestedDict(d map[string]interface{}, path []string, value interface{}) {
	if len(path) == 0 {
		return
	}
	current := d
	for _, key := range path[:len(path)-1] {
		next, ok := current[key]
		if !ok {
			next = map[string]interface{}{}
			current[key] = next
		}
		if nextDict, ok := next.(map[string]interface{}); ok {
			current = nextDict
		} else {
			newDict := map[string]interface{}{}
			current[key] = newDict
			current = newDict
		}
	}
	current[path[len(path)-1]] = value
}

func listToDictPreprocess(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for k, val := range v {
			result[k] = listToDictPreprocess(val)
		}
		return result
	case []interface{}:
		result := make(map[string]interface{}, len(v))
		for i, item := range v {
			result[fmt.Sprintf("%d", i)] = listToDictPreprocess(item)
		}
		return result
	default:
		return data
	}
}

func jsonSize(v interface{}) int {
	b, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	return len(b)
}

func formatValue(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		b, _ = json.Marshal(v)
	}
	return string(b)
}

func wrapCodeBlock(content string) string {
	return "```json\n" + content + "\n```"
}

func trimBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}

func canSplitDict(d map[string]interface{}) bool {
	if len(d) > 1 {
		return true
	}
	if len(d) == 1 {
		for _, v := range d {
			if sub, ok := v.(map[string]interface{}); ok && len(sub) > 1 {
				return true
			}
		}
	}
	return false
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	allNumeric := true
	for k := range m {
		keys = append(keys, k)
		if allNumeric {
			if _, err := strconv.Atoi(k); err != nil {
				allNumeric = false
			}
		}
	}
	if allNumeric {
		sort.Slice(keys, func(i, j int) bool {
			ni, _ := strconv.Atoi(keys[i])
			nj, _ := strconv.Atoi(keys[j])
			return ni < nj
		})
	} else {
		sort.Strings(keys)
	}
	return keys
}
