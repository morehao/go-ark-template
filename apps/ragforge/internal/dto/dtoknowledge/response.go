package dtoknowledge

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type KnowledgeCreateResp struct {
	ID uint `json:"id"`
}

type KnowledgeDetailResp struct {
	ID          uint                `json:"id"`
	KbID        uint                `json:"kbId"`
	TenantID    uint                `json:"tenantId"`
	Type        model.KnowledgeType `json:"type"`
	Title       string              `json:"title"`
	Content     string              `json:"content"`
	FileURL     string              `json:"fileUrl"`
	SourceURL   string              `json:"sourceUrl"`
	ParseStatus model.ParseStatus   `json:"parseStatus"`
	FileSize    int64               `json:"fileSize"`
	CreatorID   uint                `json:"creatorId"`
	gobject.OperatorBaseInfo
}

type KnowledgePageListItem struct {
	ID          uint                `json:"id"`
	KbID        uint                `json:"kbId"`
	TenantID    uint                `json:"tenantId"`
	Type        model.KnowledgeType `json:"type"`
	Title       string              `json:"title"`
	ParseStatus model.ParseStatus   `json:"parseStatus"`
	FileSize    int64               `json:"fileSize"`
	CreatorID   uint                `json:"creatorId"`
	gobject.OperatorBaseInfo
}

type KnowledgePageListResp struct {
	List  []KnowledgePageListItem `json:"list"`
	Total int64                   `json:"total"`
}
