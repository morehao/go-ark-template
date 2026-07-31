package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/chunker"
	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/internal/retriever"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/glog"
)

func ProcessAndEmbedContent(ctx *gin.Context, kbID, knowledgeID, tenantID uint, content string) {
	engineFactory := engine.GetGlobalFactory()
	if engineFactory == nil || engineFactory.GetEmbedding() == nil {
		glog.Errorf(ctx, "[pipeline.ProcessAndEmbedContent] engine factory or embedding provider not available")
		return
	}
	embedder := engineFactory.GetEmbedding()

	chunks := chunker.SplitText(content, chunker.DefaultConfig())
	if len(chunks) == 0 {
		glog.Errorf(ctx, "[pipeline.ProcessAndEmbedContent] no chunks generated, content length=%d", len(content))
		return
	}

	glog.Infof(ctx, "[pipeline.ProcessAndEmbedContent] splitting content into %d chunks", len(chunks))

	contents := make([]string, len(chunks))
	for i, chk := range chunks {
		contents[i] = chk.Content
	}

	embeddingResp, err := embedder.CreateEmbedding(ctx, &engine.EmbeddingRequest{
		Inputs: contents,
	})
	if err != nil {
		glog.Errorf(ctx, "[pipeline.ProcessAndEmbedContent] batch embedding fail, err:%v", err)
		return
	}
	if len(embeddingResp.Embeddings) == 0 {
		glog.Errorf(ctx, "[pipeline.ProcessAndEmbedContent] batch embedding returned empty")
		return
	}

	entries := make([]retriever.IndexEntry, 0, len(chunks))
	for i, chk := range chunks {
		if i >= len(embeddingResp.Embeddings) {
			break
		}
		entries = append(entries, retriever.IndexEntry{
			ChunkID:     0,
			KnowledgeID: knowledgeID,
			KbID:        kbID,
			TenantID:    tenantID,
			Content:     chk.Content,
			SeqID:       chk.Seq,
			Vector:      embeddingResp.Embeddings[i],
		})
	}

	if err := retriever.Get().BatchIndex(ctx, entries); err != nil {
		glog.Errorf(ctx, "[pipeline.ProcessAndEmbedContent] BatchIndex fail, err:%v", err)
		return
	}

	if err := dao.NewKnowledgeDao().UpdateMap(ctx, knowledgeID, map[string]interface{}{
		"parse_status": model.ParseStatusCompleted,
	}); err != nil {
		glog.Errorf(ctx, "[pipeline.ProcessAndEmbedContent] mark parse_status completed fail, err:%v", err)
		return
	}
	glog.Infof(ctx, "[pipeline.ProcessAndEmbedContent] completed for knowledgeID=%d, chunks=%d", knowledgeID, len(chunks))
}
