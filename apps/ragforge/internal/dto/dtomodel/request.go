package dtomodel

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type ModelCreateReq struct {
	Name      string                 `json:"name" validate:"required" label:"模型名称"`
	ModelType model.ModelType        `json:"modelType" validate:"required" label:"模型类型"`
	Provider  string                 `json:"provider" validate:"required" label:"提供商"`
	ModelName string                 `json:"modelName" validate:"required" label:"模型标识"`
	Config    map[string]interface{} `json:"config" label:"模型配置"`
}

type ModelUpdateReq struct {
	ID        uint                   `json:"id" validate:"required" label:"模型ID"`
	Name      string                 `json:"name" validate:"required" label:"模型名称"`
	ModelType model.ModelType        `json:"modelType" validate:"required" label:"模型类型"`
	Provider  string                 `json:"provider" validate:"required" label:"提供商"`
	ModelName string                 `json:"modelName" validate:"required" label:"模型标识"`
	Config    map[string]interface{} `json:"config" label:"模型配置"`
}

type ModelDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"模型ID"`
}

type ModelDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"模型ID"`
}

type ModelPageListReq struct {
	gobject.PageQuery
	Name      string          `json:"name" form:"name" label:"模型名称"`
	ModelType model.ModelType `json:"modelType" form:"modelType" label:"模型类型"`
	Provider  string          `json:"provider" form:"provider" label:"提供商"`
}

type ModelTestReq struct {
	ID uint `json:"id" validate:"required" label:"模型ID"`
}
