package dtoorg

import (
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateReq struct {
	objorg.TenantBaseInfo
	AdminInfo *objorg.TenantAdminInfo `json:"adminInfo" validate:"required" label:"管理员信息"`
}

type TenantUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	objorg.TenantBaseInfo
}

type TenantDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type TenantPageListReq struct {
	gobject.PageQuery
}

type TenantDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type DepartmentCreateReq struct {
	objorg.DepartmentBaseInfo
}
type DepartmentUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	objorg.DepartmentBaseInfo
}
type DepartmentDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}
type DepartmentPageListReq struct {
	gobject.PageQuery
}
type DepartmentDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}
