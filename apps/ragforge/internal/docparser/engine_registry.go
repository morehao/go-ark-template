package docparser

type EngineRegistration interface {
	Name() string
	Description() string
	FileTypes(docreaderConnected bool) []string
	CheckAvailable(docreaderConnected bool, overrides map[string]string) (available bool, reason string)
}

var localEngines []EngineRegistration

func RegisterEngine(e EngineRegistration) {
	localEngines = append(localEngines, e)
}

func init() {
	RegisterEngine(&builtinEngine{})
	RegisterEngine(&simpleEngine{})
}

// ---------------------------------------------------------------------------
// builtin — DocReader-backed parser for complex document formats.
// ---------------------------------------------------------------------------

type builtinEngine struct{}

func (e *builtinEngine) Name() string { return "builtin" }
func (e *builtinEngine) Description() string {
	return "DocReader built-in parser engine"
}
func (e *builtinEngine) FileTypes(_ bool) []string {
	return []string{"docx", "doc", "pdf", "md", "markdown", "xlsx", "xls", "jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp", "mp3", "wav", "m4a", "flac", "ogg"}
}
func (e *builtinEngine) CheckAvailable(docreaderConnected bool, _ map[string]string) (bool, string) {
	if docreaderConnected {
		return true, ""
	}
	return false, "DocReader service not connected"
}

const SimpleEngineName = "simple"

// ---------------------------------------------------------------------------
// simple — Go parses md/txt/csv natively
// ---------------------------------------------------------------------------

type simpleEngine struct{}

func (e *simpleEngine) Name() string { return SimpleEngineName }
func (e *simpleEngine) Description() string {
	return "Simple format & image parsing (no external service required)"
}
func (e *simpleEngine) FileTypes(_ bool) []string {
	return []string{"md", "markdown", "txt", "csv", "json", "jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp", "mp3", "wav", "m4a", "flac", "ogg"}
}
func (e *simpleEngine) CheckAvailable(_ bool, _ map[string]string) (bool, string) {
	return true, ""
}

// ---------------------------------------------------------------------------
// ListAllEngines — merge local + remote
// ---------------------------------------------------------------------------

func ListAllEngines(docreaderConnected bool, overrides map[string]string, remoteEngines []ParserEngineInfo) []ParserEngineInfo {
	remoteMap := make(map[string]ParserEngineInfo, len(remoteEngines))
	for _, re := range remoteEngines {
		remoteMap[re.Name] = re
	}

	seen := make(map[string]bool, len(localEngines))
	result := make([]ParserEngineInfo, 0, len(localEngines)+len(remoteEngines))

	for _, e := range localEngines {
		name := e.Name()
		seen[name] = true

		fileTypes := e.FileTypes(docreaderConnected)
		description := e.Description()

		if re, ok := remoteMap[name]; ok {
			if len(re.FileTypes) > 0 {
				fileTypes = re.FileTypes
			}
			if re.Description != "" {
				description = re.Description
			}
		}

		available, reason := e.CheckAvailable(docreaderConnected, overrides)
		result = append(result, ParserEngineInfo{
			Name:              name,
			Description:       description,
			FileTypes:         fileTypes,
			Available:         available,
			UnavailableReason: reason,
		})
	}

	for _, re := range remoteEngines {
		if seen[re.Name] {
			continue
		}
		result = append(result, re)
	}

	return result
}

