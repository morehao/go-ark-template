package dtochunk

import "github.com/morehao/golib/biz/gobject"

type ChunkPageListReq struct {
	gobject.PageQuery
	KnowledgeID uint `json:"knowledgeId" form:"knowledgeId" label:"知识ID"`
	KbID        uint `json:"kbId" form:"kbId" label:"知识库ID"`
}

type ChunkUpdateReq struct {
	ID      uint   `json:"id" validate:"required" label:"块ID"`
	Content string `json:"content" label:"内容"`
}

type ChunkDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"块ID"`
}

type ChunkSearchReq struct {
	KbID  uint   `json:"kbId" validate:"required" label:"知识库ID"`
	Query string `json:"query" validate:"required" label:"查询内容"`
	TopK  int    `json:"topK" label:"返回数量"`
}

type ChunkDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"块ID"`
}
