package dtopermission

import (
	"github.com/morehao/goark/apps/iam/object/objpermission"
	"github.com/morehao/golib/biz/gobject"
)

type MenuCreateReq struct {
	objpermission.MenuBaseInfo
}

type MenuUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
	objpermission.MenuBaseInfo
}

type MenuDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type MenuPageListReq struct {
	gobject.PageQuery
}

type MenuDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type RoleCreateReq struct {
	objpermission.RoleBaseInfo
}

type RoleUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
	objpermission.RoleBaseInfo
}

type RoleDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type RolePageListReq struct {
	gobject.PageQuery
}

type RoleDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type RoleAssignMenusReq struct {
	RoleID  uint   `json:"roleId" validate:"required" label:"角色ID"`    // 角色ID
	MenuIDs []uint `json:"menuIds" validate:"required" label:"菜单ID列表"` // 菜单ID列表
}

type RoleListMenusReq struct {
	RoleID uint `json:"roleId" form:"roleId" validate:"required" label:"角色ID"` // 角色ID
}

type MenuTreeReq struct {
	ParentID *uint `json:"parentID" form:"parentID" label:"父菜单ID"` // 父菜单ID
}
