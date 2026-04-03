package dtopermission

import (
	"github.com/morehao/goark/apps/iam/object/objpermission"
	"github.com/morehao/golib/biz/gobject"
)

type MenuCreateReq struct {
	objpermission.MenuBaseInfo
}

type MenuUpdateReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	objpermission.MenuBaseInfo
}

type MenuDetailReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type MenuPageListReq struct {
	gobject.PageQuery
}

type MenuDeleteReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type RoleCreateReq struct {
	objpermission.RoleBaseInfo
}
type RoleUpdateReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	objpermission.RoleBaseInfo
}
type RoleDetailReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}
type RolePageListReq struct {
	gobject.PageQuery
}
type RoleDeleteReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

// RoleAssignMenusReq 角色分配菜单请求
type RoleAssignMenusReq struct {
	// RoleID 角色ID
	RoleID uint `json:"roleId" validate:"required" label:"角色ID"`
	// MenuIDs 菜单ID列表
	MenuIDs []uint `json:"menuIds" validate:"required" label:"菜单ID列表"`
}

// RoleListMenusReq 查询角色菜单请求
type RoleListMenusReq struct {
	// RoleID 角色ID
	RoleID uint `json:"roleId" form:"roleId" validate:"required" label:"角色ID"`
}

// MenuTreeReq 菜单树请求
type MenuTreeReq struct {
	// ParentID 父菜单ID
	ParentID *uint `json:"parentID" form:"parentID" label:"父菜单ID"`
}
