package dtoapplication

import (
	"github.com/morehao/goark/apps/iam/object/objapplication"
	"github.com/morehao/golib/biz/gobject"
)

type ApplicationCreateReq struct {
	objapplication.ApplicationBaseInfo
}

type ApplicationUpdateReq struct {
	AppID uint `json:"appID" validate:"required" label:"数据自增id"` // 数据自增 ID
	objapplication.ApplicationBaseInfo
}

type ApplicationDetailReq struct {
	AppID uint `json:"appID" form:"appID" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type ApplicationPageListReq struct {
	gobject.PageQuery
}

type ApplicationDeleteReq struct {
	AppID uint `json:"appID" form:"appID" validate:"required" label:"数据自增id"` // 数据自增 ID
}
