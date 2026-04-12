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
type DepartmentTreeReq struct {
	ParentID *uint `json:"parentID" form:"parentID" label:"父部门ID"`
}
type DepartmentDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type OrganizationConfigItem struct {
	Key   string `json:"key" validate:"required" label:"配置键"`
	Value string `json:"value" label:"配置值"`
}

type OrganizationCreateReq struct {
	objorg.OrganizationBaseInfo
	Admin   *objorg.OrganizationAdminInfo `json:"admin" validate:"required" label:"管理员信息"`
	Configs []OrganizationConfigItem      `json:"configs" label:"初始化配置"`
}

type OrganizationUpdateReq struct {
	ID      uint                     `json:"id" validate:"required" label:"数据自增id"`
	Configs []OrganizationConfigItem `json:"configs" label:"配置信息"`
	objorg.OrganizationBaseInfo
}

type OrganizationDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type OrganizationGetConfigsByDomainReq struct {
	Domain string `json:"domain" form:"domain" label:"组织域名"`
}

type OrganizationPageListReq struct {
	gobject.PageQuery
	Name string `json:"name" form:"name" label:"组织名称"`
}

type OrganizationDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}
