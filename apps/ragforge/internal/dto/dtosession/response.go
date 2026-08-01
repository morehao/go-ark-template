package dtosession

import "github.com/morehao/golib/biz/gobject"

type SessionCreateResp struct {
	ID uint `json:"id"`
}

type SessionDetailResp struct {
	ID          uint   `json:"id"`
	TenantID    uint   `json:"tenantId"`
	UserID      uint   `json:"userId"`
	KbID        uint   `json:"kbId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPinned    bool   `json:"isPinned"`
	gobject.OperatorBaseInfo
}

type SessionPageListItem struct {
	ID          uint   `json:"id"`
	TenantID    uint   `json:"tenantId"`
	UserID      uint   `json:"userId"`
	KbID        uint   `json:"kbId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPinned    bool   `json:"isPinned"`
	gobject.OperatorBaseInfo
}

type SessionPageListResp struct {
	List  []SessionPageListItem `json:"list"`
	Total int64                 `json:"total"`
}

type SessionGenerateTitleResp struct {
	Title string `json:"title"`
}
