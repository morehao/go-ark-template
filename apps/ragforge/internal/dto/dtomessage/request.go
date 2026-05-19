package dtomessage

type MessageListReq struct {
	SessionID uint `json:"sessionId" form:"sessionId" validate:"required" label:"会话ID"`
}

type MessageDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"消息ID"`
}

type MessageSearchReq struct {
	SessionID uint   `json:"sessionId" validate:"required" label:"会话ID"`
	Keyword   string `json:"keyword" validate:"required" label:"搜索关键词"`
	Page      int    `json:"page" label:"页码"`
	PageSize  int    `json:"pageSize" label:"每页数量"`
}
