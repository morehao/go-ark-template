package dtokb

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type KBCreateResp struct {
	ID uint `json:"id"`
}

type KBDetailResp struct {
	ID            uint               `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	KBType        model.KBType       `json:"kbType"`
	ParserEngine  model.ParserEngine `json:"parserEngine"`
	CreatorID     uint               `json:"creatorId"`
	TenantID      uint               `json:"tenantId"`
	gobject.OperatorBaseInfo
}

type KBPageListItem struct {
	ID            uint               `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	KBType        model.KBType       `json:"kbType"`
	ParserEngine  model.ParserEngine `json:"parserEngine"`
	CreatorID     uint               `json:"creatorId"`
	TenantID      uint               `json:"tenantId"`
	gobject.OperatorBaseInfo
}

type KBPageListResp struct {
	List  []KBPageListItem `json:"list"`
	Total int64            `json:"total"`
}
