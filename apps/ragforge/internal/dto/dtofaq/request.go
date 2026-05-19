package dtofaq

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type FAQCreateReq struct {
	KbID     uint            `json:"kbId" validate:"required" label:"知识库ID"`
	Question string          `json:"question" validate:"required" label:"问题"`
	Answer   string          `json:"answer" validate:"required" label:"答案"`
	Status   model.FAQStatus `json:"status" label:"状态"`
}

type FAQUpdateReq struct {
	ID       uint            `json:"id" validate:"required" label:"FAQ ID"`
	Question string          `json:"question" validate:"required" label:"问题"`
	Answer   string          `json:"answer" validate:"required" label:"答案"`
	Status   model.FAQStatus `json:"status" label:"状态"`
}

type FAQDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"FAQ ID"`
}

type FAQDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"FAQ ID"`
}

type FAQPageListReq struct {
	gobject.PageQuery
	KbID     uint            `json:"kbId" form:"kbId" label:"知识库ID"`
	Question string          `json:"question" form:"question" label:"问题"`
	Status   model.FAQStatus `json:"status" form:"status" label:"状态"`
}

type FAQSearchReq struct {
	KbID  uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Query string `json:"query" validate:"required" label:"查询内容"`
	TopK  int    `json:"topK" label:"返回数量"`
}

type FAQImportReq struct {
	KbID     uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Question string `json:"question" validate:"required" label:"问题"`
	Answer   string `json:"answer" validate:"required" label:"答案"`
}
