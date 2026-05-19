package svcchunk

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtochunk"
	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/internal/retriever"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type ChunkSvc interface {
	PageList(ctx *gin.Context, req *dtochunk.ChunkPageListReq) (*dtochunk.ChunkPageListResp, error)
	Update(ctx *gin.Context, req *dtochunk.ChunkUpdateReq) error
	Delete(ctx *gin.Context, req *dtochunk.ChunkDeleteReq) error
	Search(ctx *gin.Context, req *dtochunk.ChunkSearchReq) (*dtochunk.ChunkSearchResp, error)
	Detail(ctx *gin.Context, req *dtochunk.ChunkDetailReq) (*dtochunk.ChunkDetailResp, error)
}

type chunkSvc struct {
}

var _ ChunkSvc = (*chunkSvc)(nil)

func NewChunkSvc() ChunkSvc {
	return &chunkSvc{}
}

func (svc *chunkSvc) PageList(ctx *gin.Context, req *dtochunk.ChunkPageListReq) (*dtochunk.ChunkPageListResp, error) {
	cond := &dao.ChunkCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		KnowledgeID: req.KnowledgeID,
		KbID:        req.KbID,
	}
	dataList, total, err := dao.NewChunkDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcchunk.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ChunkGetPageListError)
	}
	list := make([]dtochunk.ChunkPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtochunk.ChunkPageListItem{
			ID:          v.ID,
			KnowledgeID: v.KnowledgeID,
			KbID:        v.KbID,
			Content:     v.Content,
			SeqID:       v.SeqID,
			Tokens:      v.Tokens,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtochunk.ChunkPageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *chunkSvc) Update(ctx *gin.Context, req *dtochunk.ChunkUpdateReq) error {
	updateEntity := &model.ChunkEntity{
		Content: req.Content,
	}
	if err := dao.NewChunkDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcchunk.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ChunkUpdateError)
	}
	return nil
}

func (svc *chunkSvc) Delete(ctx *gin.Context, req *dtochunk.ChunkDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewChunkDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcchunk.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ChunkDeleteError)
	}
	return nil
}

func (svc *chunkSvc) Search(ctx *gin.Context, req *dtochunk.ChunkSearchReq) (*dtochunk.ChunkSearchResp, error) {
	engineFactory := engine.GetGlobalFactory()
	if engineFactory == nil || engineFactory.GetEmbedding() == nil {
		return svc.searchFallback(ctx, req)
	}

	embeddingProvider := engineFactory.GetEmbedding()
	embeddingResp, err := embeddingProvider.CreateEmbedding(ctx, &engine.EmbeddingRequest{
		Inputs: []string{req.Query},
	})
	if err != nil {
		glog.Errorf(ctx, "[svcchunk.Search] CreateEmbedding fail, err:%v, query:%s", err, req.Query)
		return svc.searchFallback(ctx, req)
	}
	if len(embeddingResp.Embeddings) == 0 {
		return &dtochunk.ChunkSearchResp{List: []dtochunk.ChunkSearchItem{}}, nil
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 10
	}

	params := retriever.SearchParams{
		KbID:      req.KbID,
		TenantID:  gincontext.GetTenantID(ctx),
		Query:     req.Query,
		TopK:      topK,
		Embedding: embeddingResp.Embeddings[0],
	}
	results, err := retriever.Get().VectorSearch(ctx, params)
	if err != nil {
		glog.Errorf(ctx, "[svcchunk.Search] VectorSearch fail, err:%v", err)
		return svc.searchFallback(ctx, req)
	}

	list := make([]dtochunk.ChunkSearchItem, 0, len(results))
	for _, v := range results {
		list = append(list, dtochunk.ChunkSearchItem{
			ID:          v.ChunkID,
			KnowledgeID: v.KnowledgeID,
			Content:     v.Content,
			Score:       v.Score,
			SeqID:       v.SeqID,
		})
	}
	return &dtochunk.ChunkSearchResp{
		List: list,
	}, nil
}

func (svc *chunkSvc) Detail(ctx *gin.Context, req *dtochunk.ChunkDetailReq) (*dtochunk.ChunkDetailResp, error) {
	detailEntity, err := dao.NewChunkDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcchunk.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ChunkGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.ChunkGetDetailError)
	}
	return &dtochunk.ChunkDetailResp{
		ID:          detailEntity.ID,
		KnowledgeID: detailEntity.KnowledgeID,
		KbID:        detailEntity.KbID,
		Content:     detailEntity.Content,
		SeqID:       detailEntity.SeqID,
		Tokens:      detailEntity.Tokens,
		MetaInfo:    detailEntity.MetaInfo,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}, nil
}

func (svc *chunkSvc) searchFallback(ctx *gin.Context, req *dtochunk.ChunkSearchReq) (*dtochunk.ChunkSearchResp, error) {
	cond := &dao.ChunkCond{
		KbID: req.KbID,
	}
	dataList, err := dao.NewChunkDao().GetListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcchunk.searchFallback] GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ChunkSearchError)
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 10
	}

	list := make([]dtochunk.ChunkSearchItem, 0)
	for _, v := range dataList {
		list = append(list, dtochunk.ChunkSearchItem{
			ID:          v.ID,
			KnowledgeID: v.KnowledgeID,
			Content:     v.Content,
			SeqID:       v.SeqID,
		})
		if len(list) >= topK {
			break
		}
	}
	return &dtochunk.ChunkSearchResp{
		List: list,
	}, nil
}
