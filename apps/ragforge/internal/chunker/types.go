package chunker

import "strings"

// Chunk represents a piece of split text with position tracking.
//
// Content holds exactly the text from the original document between Start
// and End (rune offsets), so End-Start == utf8.RuneCountInString(Content).
// This invariant is relied on by document-reconstruction code paths.
//
// ContextHeader is a separately-tracked context string (e.g. a Markdown
// heading breadcrumb) that should be prepended at embedding/retrieval time
// but is NOT part of Content.
type Chunk struct {
	Content       string
	ContextHeader string
	Seq           int
	Start         int
	End           int
}

// EmbeddingContent returns the text that should be fed to the embedding
// model — the ContextHeader prepended (when set) plus the chunk content.
func (c Chunk) EmbeddingContent() string {
	body := strings.TrimSpace(c.Content)
	if c.ContextHeader == "" {
		return body
	}
	return c.ContextHeader + "\n\n" + body
}

// ImageRef is an image reference found within a chunk's content.
type ImageRef struct {
	OriginalRef string
	AltText     string
	Start       int // offset within the chunk content
	End         int
}

// SplitterConfig configures the text splitter. Strategy and TokenLimit are
// honored by the strategy entry point; the legacy SplitText path uses only
// ChunkSize/Overlap/Separators.
type SplitterConfig struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string

	Strategy  string
	TokenLimit int
	Languages  []string
}

// DefaultChunkSize = 512 chars
const DefaultChunkSize = 512

// DefaultChunkOverlap = 80 chars
const DefaultChunkOverlap = 80

// DefaultConfig returns sensible defaults.
func DefaultConfig() SplitterConfig {
	return SplitterConfig{
		ChunkSize:    DefaultChunkSize,
		ChunkOverlap: DefaultChunkOverlap,
		Separators:   []string{"\n\n", "\n", "。"},
	}
}

// ParentChildResult holds the two-level chunking output.
type ParentChildResult struct {
	Parents  []Chunk
	Children []ChildChunk
}

// ChildChunk extends Chunk with a reference to its parent.
type ChildChunk struct {
	Chunk
	ParentIndex int // index into ParentChildResult.Parents
}
