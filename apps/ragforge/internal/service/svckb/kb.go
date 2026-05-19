package svckb

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtokb"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type KBSvc interface {
	Create(ctx *gin.Context, req *dtokb.KBCreateReq) (*dtokb.KBCreateResp, error)
	Delete(ctx *gin.Context, req *dtokb.KBDeleteReq) error
	Update(ctx *gin.Context, req *dtokb.KBUpdateReq) error
	Detail(ctx *gin.Context, req *dtokb.KBDetailReq) (*dtokb.KBDetailResp, error)
	PageList(ctx *gin.Context, req *dtokb.KBPageListReq) (*dtokb.KBPageListResp, error)
	Copy(ctx *gin.Context, req *dtokb.KBCopyReq) (*dtokb.KBCreateResp, error)
}

type kbSvc struct {
}

var _ KBSvc = (*kbSvc)(nil)

func NewKBSvc() KBSvc {
	return &kbSvc{}
}

func (svc *kbSvc) Create(ctx *gin.Context, req *dtokb.KBCreateReq) (*dtokb.KBCreateResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)
	insertEntity := &model.KnowledgeBaseEntity{
		TenantID:     tenantID,
		Name:         req.Name,
		Description:  req.Description,
		KBType:       req.KBType,
		ParserEngine: req.ParserEngine,
		CreatorID:    userID,
	}
	if err := dao.NewKnowledgeBaseDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svckb.Create] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KBCreateError)
	}
	return &dtokb.KBCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *kbSvc) Delete(ctx *gin.Context, req *dtokb.KBDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewKnowledgeBaseDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svckb.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.KBDeleteError)
	}
	return nil
}

func (svc *kbSvc) Update(ctx *gin.Context, req *dtokb.KBUpdateReq) error {
	updateEntity := &model.KnowledgeBaseEntity{
		Name:         req.Name,
		Description:  req.Description,
		KBType:       req.KBType,
		ParserEngine: req.ParserEngine,
	}
	if err := dao.NewKnowledgeBaseDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svckb.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.KBUpdateError)
	}
	return nil
}

func (svc *kbSvc) Detail(ctx *gin.Context, req *dtokb.KBDetailReq) (*dtokb.KBDetailResp, error) {
	detailEntity, err := dao.NewKnowledgeBaseDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svckb.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KBGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.KBNotExistError)
	}
	resp := &dtokb.KBDetailResp{
		ID:           detailEntity.ID,
		Name:         detailEntity.Name,
		Description:  detailEntity.Description,
		KBType:       detailEntity.KBType,
		ParserEngine: detailEntity.ParserEngine,
		CreatorID:    detailEntity.CreatorID,
		TenantID:     detailEntity.TenantID,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *kbSvc) Copy(ctx *gin.Context, req *dtokb.KBCopyReq) (*dtokb.KBCreateResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	sourceEntity, err := dao.NewKnowledgeBaseDao().GetByID(ctx, req.SourceID)
	if err != nil {
		glog.Errorf(ctx, "[svckb.Copy] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KBGetDetailError)
	}
	if sourceEntity == nil || sourceEntity.ID == 0 {
		return nil, code.GetError(code.KBNotExistError)
	}

	insertEntity := &model.KnowledgeBaseEntity{
		TenantID:        tenantID,
		Name:            req.Name,
		Description:     sourceEntity.Description,
		KBType:          sourceEntity.KBType,
		ParserEngine:    sourceEntity.ParserEngine,
		EmbeddingConfig: sourceEntity.EmbeddingConfig,
		IndexStrategy:   sourceEntity.IndexStrategy,
		CreatorID:       userID,
	}
	if err := dao.NewKnowledgeBaseDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svckb.Copy] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KBCopyError)
	}
	return &dtokb.KBCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *kbSvc) PageList(ctx *gin.Context, req *dtokb.KBPageListReq) (*dtokb.KBPageListResp, error) {
	cond := &dao.KnowledgeBaseCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Name: req.Name,
	}
	dataList, total, err := dao.NewKnowledgeBaseDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svckb.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.KBGetPageListError)
	}
	list := make([]dtokb.KBPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtokb.KBPageListItem{
			ID:           v.ID,
			Name:         v.Name,
			Description:  v.Description,
			KBType:       v.KBType,
			ParserEngine: v.ParserEngine,
			CreatorID:    v.CreatorID,
			TenantID:     v.TenantID,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtokb.KBPageListResp{
		List:  list,
		Total: total,
	}, nil
}
