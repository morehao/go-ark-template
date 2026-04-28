package dtoorg

import (
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateReq struct {
	objorg.TenantBaseInfo
	AdminInfo *objorg.TenantAdminInfo `json:"adminInfo" validate:"required" label:"管理员信息"` // 管理员信息
	AppIDs    []uint                 `json:"appIDs" label:"应用ID列表"`                    // 应用ID列表
}

type TenantUpdateReq struct {
	TenantID uint `json:"tenantID" validate:"required" label:"数据自增id"` // 数据自增id
	objorg.TenantBaseInfo
	AppIDs []uint `json:"appIDs" label:"应用ID列表"` // 应用ID列表
}

type TenantDetailReq struct {
	TenantID uint `json:"tenantID" form:"tenantID" validate:"required" label:"数据自增id"` // 数据自增id
}

type TenantPageListReq struct {
	gobject.PageQuery
}

type TenantDeleteReq struct {
	TenantID uint `json:"tenantID" form:"tenantID" validate:"required" label:"数据自增id"` // 数据自增id
}

type DepartmentCreateReq struct {
	objorg.DepartmentBaseInfo
}

type DepartmentUpdateReq struct {
	DeptID uint `json:"deptID" validate:"required" label:"数据自增id"` // 数据自增id
	objorg.DepartmentBaseInfo
}

type DepartmentDetailReq struct {
	DeptID uint `json:"deptID" form:"deptID" validate:"required" label:"数据自增id"` // 数据自增id
}

type DepartmentPageListReq struct {
	gobject.PageQuery
}

type DepartmentTreeReq struct {
	ParentID *uint `json:"parentID" form:"parentID" label:"父部门ID"` // 父部门ID
}

type DepartmentDeleteReq struct {
	DeptID uint `json:"deptID" form:"deptID" validate:"required" label:"数据自增id"` // 数据自增id
}

type OrgConfigItem struct {
	Key   string `json:"key" validate:"required" label:"配置键"` // 配置键
	Value string `json:"value" label:"配置值"`                   // 配置值
}

type OrgCreateReq struct {
	objorg.OrgBaseInfo
	Admin   *objorg.OrgAdminInfo `json:"admin" validate:"required" label:"管理员信息"` // 管理员信息
	Configs []OrgConfigItem      `json:"configs" label:"初始化配置"`                   // 初始化配置
	AppIDs  []uint               `json:"appIDs" label:"应用ID列表"`                    // 应用ID列表
}

type OrgUpdateReq struct {
	OrgID   uint            `json:"orgID" validate:"required" label:"数据自增id"` // 数据自增id
	Configs []OrgConfigItem `json:"configs" label:"配置信息"`                   // 配置信息
	AppIDs  []uint          `json:"appIDs" label:"应用ID列表"`                   // 应用ID列表
	objorg.OrgBaseInfo
}

type OrgDetailReq struct {
	OrgID uint `json:"orgID" form:"orgID" validate:"required" label:"数据自增id"` // 数据自增id
}

type GetOrganizationConfigsReq struct {
	DisplayCode string `json:"displayCode" form:"displayCode" label:"组织编码"` // 组织编码
}

type OrgPageListReq struct {
	gobject.PageQuery
	Name string `json:"name" form:"name" label:"组织名称"` // 组织名称
}

type OrgDeleteReq struct {
	OrgID uint `json:"orgID" form:"orgID" validate:"required" label:"数据自增id"` // 数据自增id
}
