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

// RoleMenuListItem 角色菜单列表项
type RoleMenuListItem struct {
	// ID 菜单ID
	ID uint `json:"id"`
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
}

// RoleMenuListResp 角色菜单列表响应
type RoleMenuListResp struct {
	// List 菜单列表
	List []RoleMenuListItem `json:"list"`
}

// MenuTreeNode 菜单树节点（实现 gtree.TreeNode[uint] 接口）
type MenuTreeNode struct {
	// ID 菜单ID
	ID uint `json:"id"`
	objpermission.MenuBaseInfo
	gobject.OperatorBaseInfo
	// Children 子菜单列表（JSON 输出）
	Children []MenuTreeNode `json:"children"`
}

func (n *MenuTreeNode) GetKey() uint {
	return n.ID
}

func (n *MenuTreeNode) GetParentKey() uint {
	return n.ParentID
}

func (n *MenuTreeNode) IsRoot() bool {
	return n.ParentID == 0
}

// MenuTreeResp 菜单树响应
type MenuTreeResp struct {
	// List 菜单树列表
	List []MenuTreeNode `json:"list"`
}
