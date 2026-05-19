package dtosession

import "github.com/morehao/golib/biz/gobject"

type SessionCreateReq struct {
	KbID    uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Title   string `json:"title" label:"会话标题"`
	Content string `json:"content" label:"会话内容"`
}

type SessionUpdateReq struct {
	ID    uint   `json:"id" validate:"required" label:"会话ID"`
	Title string `json:"title" validate:"required" label:"会话标题"`
}

type SessionDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"会话ID"`
}

type SessionDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"会话ID"`
}

type SessionPageListReq struct {
	gobject.PageQuery
	KbID uint `json:"kbId" form:"kbId" label:"知识库ID"`
}

type SessionGenerateTitleReq struct {
	SessionID uint `json:"sessionId" validate:"required" label:"会话ID"`
}

type SessionStopReq struct {
	SessionID uint `json:"sessionId" validate:"required" label:"会话ID"`
}
