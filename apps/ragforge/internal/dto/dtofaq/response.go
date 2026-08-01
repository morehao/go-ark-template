package dtofaq

import (
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gobject"
)

type FAQCreateResp struct {
	ID uint `json:"id"`
}

type FAQDetailResp struct {
	ID               uint            `json:"id"`
	KbID             uint            `json:"kbId"`
	TenantID         uint            `json:"tenantId"`
	Question         string          `json:"question"`
	Answer           string          `json:"answer"`
	SimilarQuestions string          `json:"similarQuestions"`
	Tags             string          `json:"tags"`
	Status           model.FAQStatus `json:"status"`
	CreatorID        uint            `json:"creatorId"`
	gobject.OperatorBaseInfo
}

type FAQPageListItem struct {
	ID        uint            `json:"id"`
	KbID      uint            `json:"kbId"`
	Question  string          `json:"question"`
	Answer    string          `json:"answer"`
	Status    model.FAQStatus `json:"status"`
	CreatorID uint            `json:"creatorId"`
	gobject.OperatorBaseInfo
}

type FAQPageListResp struct {
	List  []FAQPageListItem `json:"list"`
	Total int64             `json:"total"`
}

type FAQSearchItem struct {
	ID       uint    `json:"id"`
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
	Score    float64 `json:"score"`
}

type FAQSearchResp struct {
	List []FAQSearchItem `json:"list"`
}
