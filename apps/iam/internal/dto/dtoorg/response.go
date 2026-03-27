package dtoorg

import (
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateResp struct {
	ID uint `json:"id"`
}

type TenantDetailResp struct {
	ID uint `json:"id" validate:"required"`
	objorg.TenantBaseInfo
	gobject.OperatorBaseInfo
}

type TenantPageListItem struct {
	ID uint `json:"id" validate:"required"`
	objorg.TenantBaseInfo
	gobject.OperatorBaseInfo
}

type TenantPageListResp struct {
	List  []TenantPageListItem `json:"list"`
	Total int64                `json:"total"`
}

type DepartmentCreateResp struct {
	ID uint `json:"id"`
}
type DepartmentDetailResp struct {
	ID uint `json:"id" validate:"required"`
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}
type DepartmentPageListItem struct {
	ID uint `json:"id" validate:"required"`
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}
type DepartmentPageListResp struct {
	List  []DepartmentPageListItem `json:"list"`
	Total int64                    `json:"total"`
}
