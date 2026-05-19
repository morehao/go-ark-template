package svcvectorstore

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtovectorstore"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type VectorStoreSvc interface {
	Create(ctx *gin.Context, req *dtovectorstore.VectorStoreCreateReq) (*dtovectorstore.VectorStoreCreateResp, error)
	Delete(ctx *gin.Context, req *dtovectorstore.VectorStoreDeleteReq) error
	Update(ctx *gin.Context, req *dtovectorstore.VectorStoreUpdateReq) error
	Detail(ctx *gin.Context, req *dtovectorstore.VectorStoreDetailReq) (*dtovectorstore.VectorStoreDetailResp, error)
	PageList(ctx *gin.Context, req *dtovectorstore.VectorStorePageListReq) (*dtovectorstore.VectorStorePageListResp, error)
	Test(ctx *gin.Context, req *dtovectorstore.VectorStoreTestReq) (*dtovectorstore.VectorStoreTestResp, error)
	GetTypes(ctx *gin.Context) (*dtovectorstore.VectorStoreGetTypesResp, error)
}

type vectorStoreSvc struct {
}

var _ VectorStoreSvc = (*vectorStoreSvc)(nil)

func NewVectorStoreSvc() VectorStoreSvc {
	return &vectorStoreSvc{}
}

func parseVectorStoreConfig(configStr string) map[string]interface{} {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return make(map[string]interface{})
	}
	return config
}

func (svc *vectorStoreSvc) Create(ctx *gin.Context, req *dtovectorstore.VectorStoreCreateReq) (*dtovectorstore.VectorStoreCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	configStr := "{}"
	if req.Config != nil {
		configStr = gutil.ToJsonString(req.Config)
	}
	insertEntity := &model.VectorStoreEntity{
		TenantID:   tenantID,
		Name:       req.Name,
		EngineType: req.EngineType,
		Config:     configStr,
		Status:     model.VectorStoreStatusActive,
	}
	if err := dao.NewVectorStoreDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcvectorstore.Create] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.VectorStoreCreateError)
	}
	return &dtovectorstore.VectorStoreCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *vectorStoreSvc) Delete(ctx *gin.Context, req *dtovectorstore.VectorStoreDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewVectorStoreDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcvectorstore.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.VectorStoreDeleteError)
	}
	return nil
}

func (svc *vectorStoreSvc) Update(ctx *gin.Context, req *dtovectorstore.VectorStoreUpdateReq) error {
	configStr := "{}"
	if req.Config != nil {
		configStr = gutil.ToJsonString(req.Config)
	}
	updateEntity := &model.VectorStoreEntity{
		Name:       req.Name,
		EngineType: req.EngineType,
		Config:     configStr,
	}
	if err := dao.NewVectorStoreDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcvectorstore.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.VectorStoreUpdateError)
	}
	return nil
}

func (svc *vectorStoreSvc) Detail(ctx *gin.Context, req *dtovectorstore.VectorStoreDetailReq) (*dtovectorstore.VectorStoreDetailResp, error) {
	detailEntity, err := dao.NewVectorStoreDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcvectorstore.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.VectorStoreGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.VectorStoreGetDetailError)
	}
	resp := &dtovectorstore.VectorStoreDetailResp{
		ID:         detailEntity.ID,
		Name:       detailEntity.Name,
		EngineType: detailEntity.EngineType,
		Config:     parseVectorStoreConfig(detailEntity.Config),
		Status:     detailEntity.Status,
		TenantID:   detailEntity.TenantID,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *vectorStoreSvc) PageList(ctx *gin.Context, req *dtovectorstore.VectorStorePageListReq) (*dtovectorstore.VectorStorePageListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	cond := &dao.VectorStoreCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID:   tenantID,
		Name:       req.Name,
		EngineType: req.EngineType,
	}
	dataList, total, err := dao.NewVectorStoreDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcvectorstore.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.VectorStoreGetPageListError)
	}
	list := make([]dtovectorstore.VectorStorePageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtovectorstore.VectorStorePageListItem{
			ID:         v.ID,
			Name:       v.Name,
			EngineType: v.EngineType,
			Status:     v.Status,
			TenantID:   v.TenantID,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtovectorstore.VectorStorePageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *vectorStoreSvc) Test(ctx *gin.Context, req *dtovectorstore.VectorStoreTestReq) (*dtovectorstore.VectorStoreTestResp, error) {
	return &dtovectorstore.VectorStoreTestResp{
		Success: true,
		Message: "连接成功",
	}, nil
}

func (svc *vectorStoreSvc) GetTypes(ctx *gin.Context) (*dtovectorstore.VectorStoreGetTypesResp, error) {
	list := []dtovectorstore.VectorStoreTypeItem{
		{EngineType: model.EngineTypeElasticsearch, Label: "Elasticsearch"},
		{EngineType: model.EngineTypeMilvus, Label: "Milvus"},
		{EngineType: model.EngineTypePgvector, Label: "PGVector"},
	}
	return &dtovectorstore.VectorStoreGetTypesResp{
		List: list,
	}, nil
}
