package chatpipeline

import (
	"context"
	"fmt"

	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/internal/retriever"
	"github.com/morehao/golib/glog"
)

type PluginSearch struct{}

func NewPluginSearch(em *EventManager) *PluginSearch {
	p := &PluginSearch{}
	em.Register(p)
	return p
}

func (p *PluginSearch) ActivationEvents() []EventType {
	return []EventType{EventSearch}
}

func (p *PluginSearch) OnEvent(ctx context.Context, _ EventType, cm *ChatManage, next func() *PluginError) *PluginError {
	if cm.EmbeddingProvider == nil {
		glog.Errorf(ctx, "[PluginSearch] embedding provider not available")
		return next()
	}

	embeddingResp, err := cm.EmbeddingProvider.CreateEmbedding(ctx, &engine.EmbeddingRequest{
		Inputs: []string{cm.Query},
	})
	if err != nil {
		glog.Errorf(ctx, "[PluginSearch] CreateEmbedding fail, err:%v", err)
		return next()
	}
	if len(embeddingResp.Embeddings) == 0 {
		glog.Errorf(ctx, "[PluginSearch] empty embedding")
		return next()
	}

	params := retriever.SearchParams{
		KbID:      cm.KbID,
		TenantID:  cm.TenantID,
		Query:     cm.Query,
		TopK:      retriever.DefaultTopK,
		Embedding: embeddingResp.Embeddings[0],
	}
	results, err := retriever.Get().VectorSearch(ctx, params)
	if err != nil {
		glog.Errorf(ctx, "[PluginSearch] VectorSearch fail, err:%v", err)
		return next()
	}

	cm.SearchResults = make([]SearchResult, 0, len(results))
	for _, r := range results {
		cm.SearchResults = append(cm.SearchResults, SearchResult{
			ChunkID: r.ChunkID,
			Content: r.Content,
			Score:   r.Score,
		})
	}

	glog.Infof(ctx, "[PluginSearch] found %d results for query:%s", len(cm.SearchResults), cm.Query)
	return next()
}

func (p *PluginSearch) String() string {
	return fmt.Sprintf("PluginSearch")
}
