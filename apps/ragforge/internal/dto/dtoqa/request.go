package dtoqa

type KnowledgeChatReq struct {
	SessionID uint   `json:"sessionId" label:"会话ID"`
	KbID      uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Content   string `json:"content" validate:"required" label:"用户消息"`
}
