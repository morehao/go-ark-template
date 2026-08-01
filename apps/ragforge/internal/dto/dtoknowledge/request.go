package dtoknowledge

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type KnowledgeCreateFileReq struct {
	KbID  uint   `form:"kbId" validate:"required" label:"知识库ID"`
	Title string `form:"title" validate:"required" label:"标题"`
}

type KnowledgeCreateURLReq struct {
	KbID      uint   `json:"kbId" validate:"required" label:"知识库ID"`
	SourceURL string `json:"sourceUrl" validate:"required" label:"来源地址"`
	Title     string `json:"title" validate:"required" label:"标题"`
}

type KnowledgeCreateManualReq struct {
	KbID    uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Title   string `json:"title" validate:"required" label:"标题"`
	Content string `json:"content" validate:"required" label:"内容"`
}

type KnowledgeUpdateReq struct {
	ID      uint   `json:"id" validate:"required" label:"知识ID"`
	Title   string `json:"title" validate:"required" label:"标题"`
	Content string `json:"content" label:"内容"`
}

type KnowledgeDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"知识ID"`
}

type KnowledgeDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"知识ID"`
}

type KnowledgePageListReq struct {
	gobject.PageQuery
	KbID        uint                `json:"kbId" form:"kbId" label:"知识库ID"`
	Type        model.KnowledgeType `json:"type" form:"type" label:"知识类型"`
	ParseStatus model.ParseStatus   `json:"parseStatus" form:"parseStatus" label:"解析状态"`
}

type KnowledgeReparseReq struct {
	ID uint `json:"id" validate:"required" label:"知识ID"`
}

type KnowledgeDownloadReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"知识ID"`
}
