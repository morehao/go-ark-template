package dtopermission

import (
	"github.com/morehao/goark/iam/object/objpermission"
	"github.com/morehao/golib/biz/gobject"
)

type MenuCreateResp struct {
	MenuID uint `json:"menuID"` // 数据自增 ID
}

type MenuDetailResp struct {
	MenuID uint `json:"menuID" validate:"required"` // 数据自增 ID
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

type MenuPageListItem struct {
	MenuID uint `json:"menuID" validate:"required"` // 数据自增 ID
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

type MenuPageListResp struct {
	List  []MenuPageListItem `json:"list"`  // 数据列表
	Total int64              `json:"total"` // 数据总条数
}

type RoleCreateResp struct {
	RoleID uint `json:"roleID"` // 数据自增 ID
}

type RoleDetailResp struct {
	RoleID uint `json:"roleID" validate:"required"` // 数据自增 ID
	objpermission.RoleBaseInfo
	gobject.OperatorBaseInfo
}

type RolePageListItem struct {
	RoleID uint `json:"roleID" validate:"required"` // 数据自增 ID
	objpermission.RoleBaseInfo
	gobject.OperatorBaseInfo
}

type RolePageListResp struct {
	List  []RolePageListItem `json:"list"`  // 数据列表
	Total int64              `json:"total"` // 数据总条数
}

type RoleMenuListItem struct {
	MenuID uint `json:"menuID"` // 菜单ID
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

type RoleMenuListResp struct {
	List []RoleMenuListItem `json:"list"` // 菜单列表
}

type MenuTreeNode struct {
	MenuID   uint           `json:"menuID"`   // 菜单ID
	Children []MenuTreeNode `json:"children"` // 子菜单列表（JSON 输出）
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

func (n *MenuTreeNode) GetKey() uint {
	return n.MenuID
}

func (n *MenuTreeNode) GetParentKey() uint {
	return n.ParentID
}

func (n *MenuTreeNode) IsRoot() bool {
	return n.ParentID == 0
}

type MenuTreeResp struct {
	List []MenuTreeNode `json:"list"` // 菜单树列表
}
