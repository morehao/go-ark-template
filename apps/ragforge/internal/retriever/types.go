package retriever

import "context"

type SearchResult struct {
	ChunkID     uint
	KnowledgeID uint
	Content     string
	Score       float64
	SeqID       int
}

const DefaultTopK = 10

type SearchParams struct {
	KbID      uint
	TenantID  uint
	Query     string
	TopK      int
	Threshold float64
	Embedding []float32
}

type IndexEntry struct {
	ChunkID     uint
	KnowledgeID uint
	KbID        uint
	TenantID    uint
	Content     string
	SeqID       int
	Vector      []float32
}

type Retriever interface {
	VectorSearch(ctx context.Context, params SearchParams) ([]SearchResult, error)
	KeywordSearch(ctx context.Context, params SearchParams) ([]SearchResult, error)
	HybridSearch(ctx context.Context, params SearchParams) ([]SearchResult, error)
	BatchIndex(ctx context.Context, entries []IndexEntry) error
	DeleteByKnowledgeID(ctx context.Context, knowledgeID uint, tenantID uint) error
}
