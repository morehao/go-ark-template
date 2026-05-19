package dtochunk

import "github.com/morehao/golib/biz/gobject"

type ChunkPageListItem struct {
	ID          uint   `json:"id"`
	KnowledgeID uint   `json:"knowledgeId"`
	KbID        uint   `json:"kbId"`
	Content     string `json:"content"`
	SeqID       int    `json:"seqId"`
	Tokens      int    `json:"tokens"`
	gobject.OperatorBaseInfo
}

type ChunkPageListResp struct {
	List  []ChunkPageListItem `json:"list"`
	Total int64               `json:"total"`
}

type ChunkSearchItem struct {
	ID          uint    `json:"id"`
	KnowledgeID uint    `json:"knowledgeId"`
	Content     string  `json:"content"`
	Score       float64 `json:"score"`
	SeqID       int     `json:"seqId"`
}

type ChunkSearchResp struct {
	List []ChunkSearchItem `json:"list"`
}

type ChunkDetailResp struct {
	ID          uint   `json:"id"`
	KnowledgeID uint   `json:"knowledgeId"`
	KbID        uint   `json:"kbId"`
	Content     string `json:"content"`
	SeqID       int    `json:"seqId"`
	Tokens      int    `json:"tokens"`
	MetaInfo    string `json:"metaInfo"`
	gobject.OperatorBaseInfo
}
