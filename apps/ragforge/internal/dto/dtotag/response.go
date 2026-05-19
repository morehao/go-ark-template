package dtotag

import "github.com/morehao/golib/biz/gobject"

type TagCreateResp struct {
	ID uint `json:"id"`
}

type TagItem struct {
	ID     uint   `json:"id"`
	KbID   uint   `json:"kbId"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	gobject.OperatorBaseInfo
}

type TagListResp struct {
	List []TagItem `json:"list"`
}
