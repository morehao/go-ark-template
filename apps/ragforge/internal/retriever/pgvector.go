package retriever

import (
	"context"
	"fmt"
	"sort"

	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/glog"
	"github.com/pgvector/pgvector-go"
)

type pgRetriever struct {
}

type chunkWithScore struct {
	model.ChunkEntity
	Score float64 `gorm:"column:score"`
}

func (r *pgRetriever) VectorSearch(ctx context.Context, params SearchParams) ([]SearchResult, error) {
	queryVector := pgvector.NewVector(params.Embedding)
	db := dbclient.RagForgeDB(ctx)
	var results []chunkWithScore
	topK := params.TopK
	if topK <= 0 {
		topK = DefaultTopK
	}
	sql := fmt.Sprintf(
		"SELECT *, 1 - (vector <=> ?) AS score FROM %s WHERE kb_id = ? AND tenant_id = ? AND deleted_at IS NULL ORDER BY vector <=> ? LIMIT ?",
		model.TableNameChunk,
	)
	if err := db.Raw(sql, queryVector, params.KbID, params.TenantID, queryVector, topK).Scan(&results).Error; err != nil {
		glog.Errorf(ctx, "[pgRetriever.VectorSearch] fail, err:%v, params:%+v", err, params)
		return nil, err
	}
	return toSearchResults(results), nil
}

func (r *pgRetriever) KeywordSearch(ctx context.Context, params SearchParams) ([]SearchResult, error) {
	db := dbclient.RagForgeDB(ctx)
	var chunks []model.ChunkEntity
	topK := params.TopK
	if topK <= 0 {
		topK = DefaultTopK
	}
	sql := fmt.Sprintf(
		"SELECT * FROM %s WHERE kb_id = ? AND tenant_id = ? AND deleted_at IS NULL AND content ILIKE ? LIMIT ?",
		model.TableNameChunk,
	)
	likePattern := "%" + params.Query + "%"
	if err := db.Raw(sql, params.KbID, params.TenantID, likePattern, topK).Scan(&chunks).Error; err != nil {
		glog.Errorf(ctx, "[pgRetriever.KeywordSearch] fail, err:%v, params:%+v", err, params)
		return nil, err
	}
	results := make([]SearchResult, 0, len(chunks))
	for _, c := range chunks {
		results = append(results, SearchResult{
			ChunkID:     c.ID,
			KnowledgeID: c.KnowledgeID,
			Content:     c.Content,
			Score:       1.0,
			SeqID:       c.SeqID,
		})
	}
	return results, nil
}

func (r *pgRetriever) HybridSearch(ctx context.Context, params SearchParams) ([]SearchResult, error) {
	vectorResults, vecErr := r.VectorSearch(ctx, params)
	keywordResults, kwErr := r.KeywordSearch(ctx, params)
	if vecErr != nil && kwErr != nil {
		glog.Errorf(ctx, "[pgRetriever.HybridSearch] both searches fail, vecErr:%v, kwErr:%v", vecErr, kwErr)
		return nil, vecErr
	}
	if vecErr != nil {
		return keywordResults, nil
	}
	if kwErr != nil {
		return vectorResults, nil
	}
	if len(vectorResults) == 0 && len(keywordResults) == 0 {
		return []SearchResult{}, nil
	}

	seen := make(map[uint]struct{})
	rankedMap := make(map[uint]float64)

	for i, r := range vectorResults {
		rankedMap[r.ChunkID] += 1.0 / float64(i+1)
		seen[r.ChunkID] = struct{}{}
	}
	for i, r := range keywordResults {
		rankedMap[r.ChunkID] += 1.0 / float64(i+1)
		seen[r.ChunkID] = struct{}{}
	}

	merged := make([]SearchResult, 0, len(seen))
	for _, r := range vectorResults {
		if _, ok := seen[r.ChunkID]; ok {
			r.Score = rankedMap[r.ChunkID]
			merged = append(merged, r)
		}
	}
	for _, r := range keywordResults {
		if _, exists := seen[r.ChunkID]; exists {
			continue
		}
		r.Score = rankedMap[r.ChunkID]
		merged = append(merged, r)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Score > merged[j].Score
	})

	if len(merged) > params.TopK {
		merged = merged[:params.TopK]
	}
	return merged, nil
}

func (r *pgRetriever) BatchIndex(ctx context.Context, entries []IndexEntry) error {
	entities := make([]model.ChunkEntity, 0, len(entries))
	for _, entry := range entries {
		entities = append(entities, model.ChunkEntity{
			KnowledgeID: entry.KnowledgeID,
			KbID:        entry.KbID,
			TenantID:    entry.TenantID,
			Content:     entry.Content,
			SeqID:       entry.SeqID,
			Vector:      pgvector.NewVector(entry.Vector),
		})
	}
	db := dbclient.RagForgeDB(ctx)
	if err := db.Create(&entities).Error; err != nil {
		glog.Errorf(ctx, "[pgRetriever.BatchIndex] fail, err:%v, count:%d", err, len(entries))
		return err
	}
	glog.Infof(ctx, "[pgRetriever.BatchIndex] inserted %d chunks", len(entries))
	return nil
}

func (r *pgRetriever) DeleteByKnowledgeID(ctx context.Context, knowledgeID uint, tenantID uint) error {
	db := dbclient.RagForgeDB(ctx)
	result := db.Where("knowledge_id = ? AND tenant_id = ?", knowledgeID, tenantID).Delete(&model.ChunkEntity{})
	if result.Error != nil {
		glog.Errorf(ctx, "[pgRetriever.DeleteByKnowledgeID] fail, err:%v, knowledgeID:%d", result.Error, knowledgeID)
		return result.Error
	}
	glog.Infof(ctx, "[pgRetriever.DeleteByKnowledgeID] deleted %d chunks for knowledgeID:%d", result.RowsAffected, knowledgeID)
	return nil
}

func toSearchResults(results []chunkWithScore) []SearchResult {
	out := make([]SearchResult, 0, len(results))
	for _, r := range results {
		out = append(out, SearchResult{
			ChunkID:     r.ID,
			KnowledgeID: r.KnowledgeID,
			Content:     r.Content,
			Score:       r.Score,
			SeqID:       r.SeqID,
		})
	}
	return out
}

var _ Retriever = (*pgRetriever)(nil)
