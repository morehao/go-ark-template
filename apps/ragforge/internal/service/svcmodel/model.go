package svcmodel

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/dto/dtomodel"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/genericdao"
	"github.com/morehao/golib/biz/gobject"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gutil"
)

type ModelSvc interface {
	Create(ctx *gin.Context, req *dtomodel.ModelCreateReq) (*dtomodel.ModelCreateResp, error)
	Delete(ctx *gin.Context, req *dtomodel.ModelDeleteReq) error
	Update(ctx *gin.Context, req *dtomodel.ModelUpdateReq) error
	Detail(ctx *gin.Context, req *dtomodel.ModelDetailReq) (*dtomodel.ModelDetailResp, error)
	PageList(ctx *gin.Context, req *dtomodel.ModelPageListReq) (*dtomodel.ModelPageListResp, error)
	Test(ctx *gin.Context, req *dtomodel.ModelTestReq) (*dtomodel.ModelTestResp, error)
	GetProviders(ctx *gin.Context) (*dtomodel.ModelGetProvidersResp, error)
}

type modelSvc struct {
}

var _ ModelSvc = (*modelSvc)(nil)

func NewModelSvc() ModelSvc {
	return &modelSvc{}
}

func parseConfig(configStr string) map[string]interface{} {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return make(map[string]interface{})
	}
	return config
}

func (svc *modelSvc) Create(ctx *gin.Context, req *dtomodel.ModelCreateReq) (*dtomodel.ModelCreateResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	configStr := "{}"
	if req.Config != nil {
		configStr = gutil.ToJsonString(req.Config)
	}
	insertEntity := &model.ModelEntity{
		TenantID:  tenantID,
		Name:      req.Name,
		ModelType: req.ModelType,
		Provider:  req.Provider,
		ModelName: req.ModelName,
		Config:    configStr,
		Status:    model.ModelStatusActive,
	}
	if err := dao.NewModelDao().Insert(ctx, insertEntity); err != nil {
		glog.Errorf(ctx, "[svcmodel.Create] dao Insert fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ModelCreateError)
	}
	return &dtomodel.ModelCreateResp{
		ID: insertEntity.ID,
	}, nil
}

func (svc *modelSvc) Delete(ctx *gin.Context, req *dtomodel.ModelDeleteReq) error {
	userID := gincontext.GetUserID(ctx)
	if err := dao.NewModelDao().Delete(ctx, req.ID, userID); err != nil {
		glog.Errorf(ctx, "[svcmodel.Delete] dao Delete fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ModelDeleteError)
	}
	return nil
}

func (svc *modelSvc) Update(ctx *gin.Context, req *dtomodel.ModelUpdateReq) error {
	configStr := "{}"
	if req.Config != nil {
		configStr = gutil.ToJsonString(req.Config)
	}
	updateEntity := &model.ModelEntity{
		Name:      req.Name,
		ModelType: req.ModelType,
		Provider:  req.Provider,
		ModelName: req.ModelName,
		Config:    configStr,
	}
	if err := dao.NewModelDao().UpdateByID(ctx, req.ID, updateEntity); err != nil {
		glog.Errorf(ctx, "[svcmodel.Update] dao UpdateByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return code.GetError(code.ModelUpdateError)
	}
	return nil
}

func (svc *modelSvc) Detail(ctx *gin.Context, req *dtomodel.ModelDetailReq) (*dtomodel.ModelDetailResp, error) {
	detailEntity, err := dao.NewModelDao().GetByID(ctx, req.ID)
	if err != nil {
		glog.Errorf(ctx, "[svcmodel.Detail] dao GetByID fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ModelGetDetailError)
	}
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, code.GetError(code.ModelGetDetailError)
	}
	resp := &dtomodel.ModelDetailResp{
		ID:        detailEntity.ID,
		Name:      detailEntity.Name,
		ModelType: detailEntity.ModelType,
		Provider:  detailEntity.Provider,
		ModelName: detailEntity.ModelName,
		Config:    parseConfig(detailEntity.Config),
		Status:    detailEntity.Status,
		TenantID:  detailEntity.TenantID,
		OperatorBaseInfo: gobject.OperatorBaseInfo{
			CreatedAt: detailEntity.CreatedAt.Unix(),
			UpdatedAt: detailEntity.UpdatedAt.Unix(),
		},
	}
	return resp, nil
}

func (svc *modelSvc) PageList(ctx *gin.Context, req *dtomodel.ModelPageListReq) (*dtomodel.ModelPageListResp, error) {
	tenantID := gincontext.GetTenantID(ctx)
	cond := &dao.ModelCond{
		BaseCond: &genericdao.BaseCond{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		TenantID:  tenantID,
		Name:      req.Name,
		ModelType: req.ModelType,
		Provider:  req.Provider,
	}
	dataList, total, err := dao.NewModelDao().GetPageListByCond(ctx, cond)
	if err != nil {
		glog.Errorf(ctx, "[svcmodel.PageList] dao GetPageListByCond fail, err:%v, req:%s", err, gutil.ToJsonString(req))
		return nil, code.GetError(code.ModelGetPageListError)
	}
	list := make([]dtomodel.ModelPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtomodel.ModelPageListItem{
			ID:        v.ID,
			Name:      v.Name,
			ModelType: v.ModelType,
			Provider:  v.Provider,
			ModelName: v.ModelName,
			Status:    v.Status,
			TenantID:  v.TenantID,
			OperatorBaseInfo: gobject.OperatorBaseInfo{
				UpdatedAt: v.UpdatedAt.Unix(),
			},
		})
	}
	return &dtomodel.ModelPageListResp{
		List:  list,
		Total: total,
	}, nil
}

func (svc *modelSvc) Test(ctx *gin.Context, req *dtomodel.ModelTestReq) (*dtomodel.ModelTestResp, error) {
	return &dtomodel.ModelTestResp{
		Success: true,
		Message: "连接成功",
	}, nil
}

func (svc *modelSvc) GetProviders(ctx *gin.Context) (*dtomodel.ModelGetProvidersResp, error) {
	list := []dtomodel.ModelProviderItem{
		{Provider: "openai", Label: "OpenAI"},
		{Provider: "azure", Label: "Azure OpenAI"},
		{Provider: "ollama", Label: "Ollama"},
		{Provider: "claude", Label: "Anthropic Claude"},
		{Provider: "gemini", Label: "Google Gemini"},
		{Provider: "deepseek", Label: "DeepSeek"},
		{Provider: "qwen", Label: "通义千问"},
		{Provider: "baidu", Label: "百度文心"},
	}
	return &dtomodel.ModelGetProvidersResp{
		List: list,
	}, nil
}
