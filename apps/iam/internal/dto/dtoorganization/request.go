package dtoorganization

import (
	"github.com/morehao/goark/apps/iam/object/objorganization"
	"github.com/morehao/golib/biz/gobject"
)

type OrganizationCreateReq struct {
	objorganization.OrganizationBaseInfo
}

type OrganizationUpdateReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	objorganization.OrganizationBaseInfo
}

type OrganizationDetailReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type OrganizationPageListReq struct {
	gobject.PageQuery
}

type OrganizationDeleteReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}
