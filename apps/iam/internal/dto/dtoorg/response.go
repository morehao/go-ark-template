package dtoorg

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/apps/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateResp struct {
	ID       uint `json:"id"`       // 租户ID
	AdminID  uint `json:"adminID"`  // 管理员ID
	PersonID uint `json:"personID"` // 自然人ID
}

type TenantDetailResp struct {
	ID uint `json:"id" validate:"required"` // 数据自增id
	objorg.TenantBaseInfo
	gobject.OperatorBaseInfo
}

type TenantPageListItem struct {
	ID uint `json:"id" validate:"required"` // 数据自增id
	objorg.TenantBaseInfo
	gobject.OperatorBaseInfo
}

type TenantPageListResp struct {
	List  []TenantPageListItem `json:"list"`  // 数据列表
	Total int64                `json:"total"` // 数据总条数
}

type DepartmentCreateResp struct {
	ID uint `json:"id"` // 数据自增id
}

type DepartmentDetailResp struct {
	ID uint `json:"id" validate:"required"` // 数据自增id
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}

type DepartmentPageListItem struct {
	ID uint `json:"id" validate:"required"` // 数据自增id
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}

type DepartmentPageListResp struct {
	List  []DepartmentPageListItem `json:"list"`  // 数据列表
	Total int64                    `json:"total"` // 数据总条数
}

type DepartmentTreeNode struct {
	ID       uint                 `json:"id"`       // 数据自增id
	Children []DepartmentTreeNode `json:"children"` // 子部门列表
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}

type DepartmentTreeResp struct {
	List []DepartmentTreeNode `json:"list"` // 部门树列表
}

type OrgCreateResp struct {
	ID      uint `json:"id"`      // 组织ID
	AdminID uint `json:"adminID"` // 管理员ID
}

type GetOrganizationConfigsResp struct {
	OrgID   uint                         `json:"orgId"`   // 组织ID
	OrgName string                       `json:"orgName"` // 组织名称
	Domain  string                       `json:"domain"`  // 域名
	Logo    string                       `json:"logo"`    // Logo
	Status  model.OrgStatus              `json:"status"`  // 状态
	Configs map[string]map[string]string `json:"configs"` // 配置信息
}

type OrgDetailResp struct {
	ID      uint                         `json:"id" validate:"required"` // 数据自增id
	Configs map[string]map[string]string `json:"configs"`                // 配置信息
	objorg.OrgBaseInfo
	gobject.OperatorBaseInfo
}

type OrgPageListItem struct {
	ID uint `json:"id" validate:"required"` // 数据自增id
	objorg.OrgBaseInfo
	gobject.OperatorBaseInfo
}

type OrgPageListResp struct {
	List  []OrgPageListItem `json:"list"`  // 数据列表
	Total int64             `json:"total"` // 数据总条数
}

type OrgConfigOptionResp struct {
	Value       string `json:"value"`       // 配置值
	Description string `json:"description"` // 配置说明
}

type OrgConfigMetaResp struct {
	Group        string                `json:"group"`             // 配置分组
	Key          string                `json:"key"`               // 配置键
	Type         string                `json:"type"`              // 配置类型
	DefaultValue string                `json:"defaultValue"`      // 默认值
	Description  string                `json:"description"`       // 配置说明
	Options      []OrgConfigOptionResp `json:"options,omitempty"` // 配置选项
}

type OrgConfigListResp struct {
	Configs []OrgConfigMetaResp `json:"configs"` // 配置列表
}
