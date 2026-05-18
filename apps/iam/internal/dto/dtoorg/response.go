package dtoorg

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/iam/object/objorg"
	"github.com/morehao/golib/biz/gobject"
)

type TenantCreateResp struct {
	TenantID uint `json:"tenantID"` // 租户ID
	AdminID  uint `json:"adminID"`  // 管理员ID
	PersonID uint `json:"personID"` // 自然人ID
}

type TenantDetailResp struct {
	TenantID uint `json:"tenantID" validate:"required"` // 数据自增id
	Apps     []AppInfo `json:"apps"`                   // 应用列表
	objorg.TenantBaseInfo
	gobject.OperatorBaseInfo
}

type TenantPageListItem struct {
	TenantID uint `json:"tenantID" validate:"required"` // 数据自增id
	objorg.TenantBaseInfo
	gobject.OperatorBaseInfo
}

type TenantPageListResp struct {
	List  []TenantPageListItem `json:"list"`  // 数据列表
	Total int64                `json:"total"` // 数据总条数
}

type DepartmentCreateResp struct {
	DeptID uint `json:"deptID"` // 数据自增id
}

type DepartmentDetailResp struct {
	DeptID uint `json:"deptID" validate:"required"` // 数据自增id
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}

type DepartmentPageListItem struct {
	DeptID uint `json:"deptID" validate:"required"` // 数据自增id
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}

type DepartmentPageListResp struct {
	List  []DepartmentPageListItem `json:"list"`  // 数据列表
	Total int64                    `json:"total"` // 数据总条数
}

type DepartmentTreeNode struct {
	DeptID   uint                 `json:"deptID"`   // 数据自增id
	Children []DepartmentTreeNode `json:"children"` // 子部门列表
	objorg.DepartmentBaseInfo
	gobject.OperatorBaseInfo
}

type DepartmentTreeResp struct {
	List []DepartmentTreeNode `json:"list"` // 部门树列表
}

type OrgCreateResp struct {
	OrgID   uint `json:"orgID"`   // 组织ID
	AdminID uint `json:"adminID"` // 管理员ID
}

type ConfigItemResp struct {
	Key   string `json:"key"`   // 配置键
	Value string `json:"value"` // 配置值
}

type ConfigGroupResp struct {
	Group string           `json:"group"` // 配置分组
	Items []ConfigItemResp `json:"items"` // 配置项列表
}

type GetOrgConfigResp struct {
	OrgID   uint              `json:"orgId"`   // 组织ID
	OrgName string            `json:"orgName"` // 组织名称
	Logo    string            `json:"logo"`    // Logo
	Status  model.OrgStatus   `json:"status"`  // 状态
	Configs []ConfigGroupResp `json:"configs"` // 配置信息
}

type OrgDetailResp struct {
	OrgID   uint                         `json:"orgID" validate:"required"` // 数据自增id
	Configs map[string]map[string]string `json:"configs"`                   // 配置信息
	Apps    []AppInfo                    `json:"apps"`                      // 应用列表
	objorg.OrgBaseInfo
	gobject.OperatorBaseInfo
}

type AppInfo struct {
	AppID   uint   `json:"appID"`
	AppName string `json:"appName"`
}

type OrgPageListItem struct {
	OrgID uint `json:"orgID" validate:"required"` // 数据自增id
	objorg.OrgBaseInfo
	gobject.OperatorBaseInfo
}

type OrgPageListResp struct {
	List  []OrgPageListItem `json:"list"`  // 数据列表
	Total int64             `json:"total"` // 数据总条数
}

type OrgConfigOptionsItem struct {
	Value       string `json:"value"`       // 配置值
	Description string `json:"description"` // 配置说明
}

type OrgConfigMetaResp struct {
	Group        string                 `json:"group"`        // 配置分组
	Key          string                 `json:"key"`          // 配置键
	Type         string                 `json:"type"`         // 配置类型
	DefaultValue string                 `json:"defaultValue"` // 默认值
	Description  string                 `json:"description"`  // 配置说明
	Options      []OrgConfigOptionsItem `json:"options"`      // 配置选项
}

type ListConfigDefinitionsResp struct {
	Configs []OrgConfigMetaResp `json:"configs"` // 配置项元数据列表
}
