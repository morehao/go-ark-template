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

type RoleAssignMenusReq struct {
	RoleID  uint   `json:"roleId" validate:"required" label:"角色ID"`
	MenuIDs []uint `json:"menuIds" validate:"required" label:"菜单ID列表"`
}

type RoleRemoveMenusReq struct {
	RoleID  uint   `json:"roleId" validate:"required" label:"角色ID"`
	MenuIDs []uint `json:"menuIds" validate:"required" label:"菜单ID列表"`
}

type RoleMenuListReq struct {
	RoleID uint `json:"roleId" form:"roleId" validate:"required" label:"角色ID"`
}

type UserAssignRolesReq struct {
	UserID  uint   `json:"userId" validate:"required" label:"用户ID"`
	RoleIDs []uint `json:"roleIds" validate:"required" label:"角色ID列表"`
}

type UserRemoveRolesReq struct {
	UserID  uint   `json:"userId" validate:"required" label:"用户ID"`
	RoleIDs []uint `json:"roleIds" validate:"required" label:"角色ID列表"`
}

type UserRoleListReq struct {
	UserID uint `json:"userId" form:"userId" validate:"required" label:"用户ID"`
}

type UserPermissionsReq struct{}

type MenuTreeReq struct{}
