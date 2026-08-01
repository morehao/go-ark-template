package dtomodel

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type ModelCreateResp struct {
	ID uint `json:"id"`
}

type ModelDetailResp struct {
	ID        uint                   `json:"id"`
	Name      string                 `json:"name"`
	ModelType model.ModelType        `json:"modelType"`
	Provider  string                 `json:"provider"`
	ModelName string                 `json:"modelName"`
	Config    map[string]interface{} `json:"config"`
	Status    model.ModelStatus      `json:"status"`
	TenantID  uint                   `json:"tenantId"`
	gobject.OperatorBaseInfo
}

type ModelPageListItem struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	ModelType model.ModelType   `json:"modelType"`
	Provider  string            `json:"provider"`
	ModelName string            `json:"modelName"`
	Status    model.ModelStatus `json:"status"`
	TenantID  uint              `json:"tenantId"`
	gobject.OperatorBaseInfo
}

type ModelPageListResp struct {
	List  []ModelPageListItem `json:"list"`
	Total int64               `json:"total"`
}

type ModelTestResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ModelProviderItem struct {
	Provider string `json:"provider"`
	Label    string `json:"label"`
}

type ModelGetProvidersResp struct {
	List []ModelProviderItem `json:"list"`
}
