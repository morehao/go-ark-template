package dtoorg

import (
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateResp struct {
	ID       uint `json:"id"`
	AdminID  uint `json:"adminID"`
	PersonID uint `json:"personID"`
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

type DepartmentTreeNode struct {
	ID uint `json:"id"`
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
	Children []DepartmentTreeNode `json:"children"`
}

type DepartmentTreeResp struct {
	List []DepartmentTreeNode `json:"list"`
}

type OrganizationCreateResp struct {
	ID      uint `json:"id"`
	AdminID uint `json:"adminID"`
}
type OrganizationDetailResp struct {
	ID uint `json:"id" validate:"required"`
	objorg.OrganizationBaseInfo
	gobject.OperatorBaseInfo
}
type OrganizationPageListItem struct {
	ID uint `json:"id" validate:"required"`
	objorg.OrganizationBaseInfo
	gobject.OperatorBaseInfo
}
type OrganizationPageListResp struct {
	List  []OrganizationPageListItem `json:"list"`
	Total int64                      `json:"total"`
}
