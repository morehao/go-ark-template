package dtoapplication

import (
	"github.com/morehao/goark/apps/iam/object/objapplication"
	"github.com/morehao/golib/biz/gobject"
)

type ApplicationCreateResp struct {
	ID uint `json:"id"` // 数据自增 ID
}

type ApplicationDetailResp struct {
	ID uint `json:"id" validate:"required"` // 数据自增 ID
	objapplication.ApplicationBaseInfo
	gobject.OperatorBaseInfo
}

type ApplicationPageListItem struct {
	ID uint `json:"id" validate:"required"` // 数据自增 ID
	objapplication.ApplicationBaseInfo
	gobject.OperatorBaseInfo
}

type ApplicationPageListResp struct {
	List  []ApplicationPageListItem `json:"list"`  // 数据列表
	Total int64                     `json:"total"` // 数据总条数
}
