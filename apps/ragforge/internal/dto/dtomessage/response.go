package dtomessage

import "github.com/morehao/goark/ragforge/model"

type MessageListItem struct {
	ID         uint             `json:"id"`
	Role       model.MessageRole `json:"role"`
	Content    string           `json:"content"`
	Metadata   string           `json:"metadata"`
	TokenCount int              `json:"tokenCount"`
	CreatedAt  int64            `json:"createdAt"`
}

type MessageListResp struct {
	List []MessageListItem `json:"list"`
}

type MessageSearchItem struct {
	ID         uint             `json:"id"`
	SessionID  uint             `json:"sessionId"`
	Role       model.MessageRole `json:"role"`
	Content    string           `json:"content"`
	Metadata   string           `json:"metadata"`
	TokenCount int              `json:"tokenCount"`
	CreatedAt  int64            `json:"createdAt"`
}

type MessageSearchResp struct {
	List     []MessageSearchItem `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}
