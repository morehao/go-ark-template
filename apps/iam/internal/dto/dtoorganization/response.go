package dtoorganization

import (
	"github.com/morehao/goark/apps/iam/object/objorganization"
	"github.com/morehao/golib/biz/gobject"
)

type OrganizationCreateResp struct {
	// ID 数据自增 ID
	ID uint `json:"id"`
}

type OrganizationDetailResp struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objorganization.OrganizationBaseInfo
	gobject.OperatorBaseInfo
}

type OrganizationPageListItem struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objorganization.OrganizationBaseInfo
	gobject.OperatorBaseInfo
}

type OrganizationPageListResp struct {
	// List 数据列表
	List []OrganizationPageListItem `json:"list"`
	// Total 数据总条数
	Total int64 `json:"total"`
}
