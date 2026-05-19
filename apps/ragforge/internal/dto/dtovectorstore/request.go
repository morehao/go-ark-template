package dtovectorstore

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type VectorStoreCreateReq struct {
	Name       string                 `json:"name" validate:"required" label:"向量库名称"`
	EngineType model.EngineType       `json:"engineType" validate:"required" label:"引擎类型"`
	Config     map[string]interface{} `json:"config" label:"连接配置"`
}

type VectorStoreUpdateReq struct {
	ID         uint                   `json:"id" validate:"required" label:"向量库ID"`
	Name       string                 `json:"name" validate:"required" label:"向量库名称"`
	EngineType model.EngineType       `json:"engineType" validate:"required" label:"引擎类型"`
	Config     map[string]interface{} `json:"config" label:"连接配置"`
}

type VectorStoreDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"向量库ID"`
}

type VectorStoreDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"向量库ID"`
}

type VectorStorePageListReq struct {
	gobject.PageQuery
	Name       string            `json:"name" form:"name" label:"向量库名称"`
	EngineType model.EngineType  `json:"engineType" form:"engineType" label:"引擎类型"`
}

type VectorStoreTestReq struct {
	ID uint `json:"id" validate:"required" label:"向量库ID"`
}
