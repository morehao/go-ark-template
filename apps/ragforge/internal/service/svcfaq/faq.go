package svcfaq

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtofaq"
	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/internal/retriever"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/dbaccess/gormdao"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type FAQSvc interface {
	Create(ctx *gin.Context, req *dtofaq.FAQCreateReq) (*dtofaq.FAQCreateResp, error)
	Update(ctx *gin.Context, req *dtofaq.FAQUpdateReq) error
	Delete(ctx *gin.Context, req *dtofaq.FAQDeleteReq) error
	Detail(ctx *gin.Context, req *dtofaq.FAQDetailReq) (*dtofaq.FAQDetailResp, error)
	PageList(ctx *gin.Context, req *dtofaq.FAQPageListReq) (*dtofaq.FAQPageListResp, error)
	Search(ctx *gin.Context, req *dtofaq.FAQSearchReq) (*dtofaq.FAQSearchResp, error)
}

type faqSvc struct {
}

var _ FAQSvc = (*faqSvc)(nil)

func NewFAQSvc() FAQSvc {
	return &faqSvc{}
}

func (svc *faqSvc) Create(ctx *gin.Context, req *dtofaq.FAQCreateReq) (*dtofaq.FAQCreateResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)
	status := req.Status
	if status == "" {
		status = model.FAQStatusActive
	}
	insertEntity := &model.FAQEntity{
		KbID:      req.KbID,
		TenantID:  tenantID,
		Question:  req.Question,
		Answer:    req.Answer,
		Status:    status,
		CreatorID: userID,
	}
	if err := dao.NewFAQDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcfaq.Create] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.FAQCreateError)
	}
	return &dtofaq.FAQCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *faqSvc) Update(ctx *gin.Context, req *dtofaq.FAQUpdateReq) error {
	status := req.Status
	if status == "" {
		status = model.FAQStatusActive
	}
	updateEntity := &model.FAQEntity{
		Question: req.Question,
		Answer:   req.Answer,
		Status:   status,
	}
	if err := dao.NewFAQDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcfaq.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.FAQUpdateError)
	}
	return nil
}

func (svc *faqSvc) Delete(ctx *gin.Context, req *dtofaq.FAQDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewFAQDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcfaq.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.FAQDeleteError)
	}
	return nil
}

func (svc *faqSvc) Detail(ctx *gin.Context, req *dtofaq.FAQDetailReq) (*dtofaq.FAQDetailResp, error) {
	detailEntity, err := dao.NewFAQDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcfaq.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.FAQGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.FAQGetDetailError)
	}
	resp := &dtofaq.FAQDetailResp{
		ID:               detailEntity.ID,
		KbID:             detailEntity.KbID,
		TenantID:         detailEntity.TenantID,
		Question:         detailEntity.Question,
		Answer:           detailEntity.Answer,
		SimilarQuestions: detailEntity.SimilarQuestions,
		Tags:             detailEntity.Tags,
		Status:           detailEntity.Status,
		CreatorID:        detailEntity.CreatorID,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *faqSvc) PageList(ctx *gin.Context, req *dtofaq.FAQPageListReq) (*dtofaq.FAQPageListResp, error) {
	cond := &dao.FAQCond{
		BaseCond: &gormdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		KbID:     req.KbID,
		Question: req.Question,
		Status:   req.Status,
	}
	dataList, total, err := dao.NewFAQDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcfaq.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.FAQGetPageListError)
	}
	list := make([]dtofaq.FAQPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtofaq.FAQPageListItem{
			ID:        v.ID,
			KbID:      v.KbID,
			Question:  v.Question,
			Answer:    v.Answer,
			Status:    v.Status,
			CreatorID: v.CreatorID,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtofaq.FAQPageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *faqSvc) Search(ctx *gin.Context, req *dtofaq.FAQSearchReq) (*dtofaq.FAQSearchResp, error) {
	topK := req.TopK
	if topK <= 0 {
		topK = 10
	}

	tenantID := gincontext.GetTenantID(ctx)

	if r := retriever.Get(); r != nil {
		if embeddingProvider := engine.GetGlobalFactory().GetEmbedding(); embeddingProvider != nil {
			embeddingResp, err := embeddingProvider.CreateEmbedding(ctx, &engine.EmbeddingRequest{
				Model:  "",
				Inputs: []string{req.Query},
			})
			if err == nil && len(embeddingResp.Embeddings) > 0 {
				searchResults, err := r.VectorSearch(ctx, retriever.SearchParams{
					KbID:      req.KbID,
					TenantID:  tenantID,
					Query:     req.Query,
					TopK:      topK,
					Embedding: embeddingResp.Embeddings[0],
				})
				if err == nil {
					list := make([]dtofaq.FAQSearchItem, 0, len(searchResults))
					for _, sr := range searchResults {
						list = append(list, dtofaq.FAQSearchItem{
							ID:       sr.ChunkID,
							Question: "",
							Answer:   sr.Content,
							Score:    sr.Score,
						})
					}
					return &dtofaq.FAQSearchResp{
						List: list,
					}, nil
				}
				glog.Errorf(ctx, "[svcfaq.Search] VectorSearch fail, fallback to keyword, err:%v", err)
			} else {
				glog.Errorf(ctx, "[svcfaq.Search] CreateEmbedding fail, fallback to keyword, err:%v", err)
			}
		}
	}

	cond := &dao.FAQCond{
		KbID: req.KbID,
	}
	dataList, err := dao.NewFAQDao().GetListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcfaq.Search] dao GetListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.FAQSearchError)
	}

	query := strings.ToLower(req.Query)
	list := make([]dtofaq.FAQSearchItem, 0)
	for _, v := range dataList {
		if strings.Contains(strings.ToLower(v.Question), query) || strings.Contains(strings.ToLower(v.Answer), query) {
			list = append(list, dtofaq.FAQSearchItem{
				ID:       v.ID,
				Question: v.Question,
				Answer:   v.Answer,
			})
			if len(list) >= topK {
				break
			}
		}
	}
	return &dtofaq.FAQSearchResp{
		List: list,
	}, nil
}
