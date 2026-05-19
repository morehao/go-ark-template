package dtokb

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type KBCreateReq struct {
	Name         string             `json:"name" validate:"required" label:"知识库名称"`
	Description  string             `json:"description" label:"知识库描述"`
	KBType       model.KBType       `json:"kbType" validate:"required" label:"知识库类型"`
	ParserEngine model.ParserEngine `json:"parserEngine" label:"解析引擎"`
}

type KBUpdateReq struct {
	ID           uint               `json:"id" validate:"required" label:"知识库ID"`
	Name         string             `json:"name" validate:"required" label:"知识库名称"`
	Description  string             `json:"description" label:"知识库描述"`
	KBType       model.KBType       `json:"kbType" validate:"required" label:"知识库类型"`
	ParserEngine model.ParserEngine `json:"parserEngine" label:"解析引擎"`
}

type KBDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"知识库ID"`
}

type KBPageListReq struct {
	gobject.PageQuery
	Name string `json:"name" form:"name" label:"知识库名称"`
}

type KBDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"知识库ID"`
}

type KBCopyReq struct {
	SourceID uint   `json:"sourceId" validate:"required" label:"源知识库ID"`
	Name     string `json:"name" validate:"required" label:"新知识库名称"`
}
