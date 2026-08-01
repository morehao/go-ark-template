package dtovectorstore

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type VectorStoreCreateResp struct {
	ID uint `json:"id"`
}

type VectorStoreDetailResp struct {
	ID         uint                    `json:"id"`
	Name       string                  `json:"name"`
	EngineType model.EngineType        `json:"engineType"`
	Config     map[string]interface{}  `json:"config"`
	Status     model.VectorStoreStatus `json:"status"`
	TenantID   uint                    `json:"tenantId"`
	gobject.OperatorBaseInfo
}

type VectorStorePageListItem struct {
	ID         uint                    `json:"id"`
	Name       string                  `json:"name"`
	EngineType model.EngineType        `json:"engineType"`
	Status     model.VectorStoreStatus `json:"status"`
	TenantID   uint                    `json:"tenantId"`
	gobject.OperatorBaseInfo
}

type VectorStorePageListResp struct {
	List  []VectorStorePageListItem `json:"list"`
	Total int64                     `json:"total"`
}

type VectorStoreTestResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VectorStoreTypeItem struct {
	EngineType model.EngineType `json:"engineType"`
	Label      string           `json:"label"`
}

type VectorStoreGetTypesResp struct {
	List []VectorStoreTypeItem `json:"list"`
}
