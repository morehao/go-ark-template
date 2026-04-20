package dtoorg

import (
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateReq struct {
	objorg.TenantBaseInfo
	AdminInfo *objorg.TenantAdminInfo `json:"adminInfo" validate:"required" label:"管理员信息"` // 管理员信息
}

type TenantUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	objorg.TenantBaseInfo
}

type TenantDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type TenantPageListReq struct {
	gobject.PageQuery
}

type TenantDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type DepartmentCreateReq struct {
	objorg.DepartmentBaseInfo
}

type DepartmentUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	objorg.DepartmentBaseInfo
}

type DepartmentDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type DepartmentPageListReq struct {
	gobject.PageQuery
}

type DepartmentTreeReq struct {
	ParentID *uint `json:"parentID" form:"parentID" label:"父部门ID"` // 父部门ID
}

type DepartmentDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type OrgConfigItem struct {
	Key   string `json:"key" validate:"required" label:"配置键"` // 配置键
	Value string `json:"value" label:"配置值"`                   // 配置值
}

type OrgCreateReq struct {
	objorg.OrgBaseInfo
	Admin   *objorg.OrgAdminInfo  `json:"admin" validate:"required" label:"管理员信息"` // 管理员信息
	Configs []OrgConfigItem       `json:"configs" label:"初始化配置"`                   // 初始化配置
}

type OrgUpdateReq struct {
	ID      uint              `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	Configs []OrgConfigItem   `json:"configs" label:"配置信息"`                  // 配置信息
	objorg.OrgBaseInfo
}

type OrgDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type OrgGetConfigsByDomainReq struct {
	Domain string `json:"domain" form:"domain" label:"组织域名"` // 组织域名
}

type OrgPageListReq struct {
	gobject.PageQuery
	Name string `json:"name" form:"name" label:"组织名称"` // 组织名称
}

type OrgDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}