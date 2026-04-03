package dtopermission

import (
	"github.com/morehao/goark/apps/iam/object/objpermission"
	"github.com/morehao/golib/biz/gobject"
)

type MenuCreateResp struct {
	// ID 数据自增 ID
	ID uint `json:"id"`
}

type MenuDetailResp struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

type MenuPageListItem struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

type MenuPageListResp struct {
	// List 数据列表
	List []MenuPageListItem `json:"list"`
	// Total 数据总条数
	Total int64 `json:"total"`
}

type RoleCreateResp struct {
	// ID 数据自增 ID
	ID uint `json:"id"`
}
type RoleDetailResp struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objpermission.RoleBaseInfo
	gobject.OperatorBaseInfo
}
type RolePageListItem struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objpermission.RoleBaseInfo
	gobject.OperatorBaseInfo
}
type RolePageListResp struct {
	// List 数据列表
	List []RolePageListItem `json:"list"`
	// Total 数据总条数
	Total int64 `json:"total"`
}

type RoleMenuListResp struct {
	List []RoleMenuItem `json:"list"`
}

type RoleMenuItem struct {
	MenuID     uint   `json:"menuId"`
	MenuName   string `json:"menuName"`
	MenuCode   string `json:"menuCode"`
	MenuType   string `json:"menuType"`
	Permission string `json:"permission"`
	ParentID   uint   `json:"parentId"`
}

type UserRoleListResp struct {
	List []UserRoleItem `json:"list"`
}

type UserRoleItem struct {
	RoleID   uint   `json:"roleId"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
	RoleType string `json:"roleType"`
}

type UserPermissionsResp struct {
	Menus       []MenuTreeNode `json:"menus"`
	Permissions []string       `json:"permissions"`
}

type MenuTreeNode struct {
	ID            uint           `json:"id"`
	MenuName      string         `json:"menuName"`
	MenuCode      string         `json:"menuCode"`
	MenuType      string         `json:"menuType"`
	ParentID      uint           `json:"parentId"`
	RoutePath     string         `json:"routePath"`
	ComponentPath string         `json:"componentPath"`
	Permission    string         `json:"permission"`
	Icon          string         `json:"icon"`
	SortOrder     int32          `json:"sortOrder"`
	Visibility    string         `json:"visibility"`
	Children      []MenuTreeNode `json:"children"`
}
