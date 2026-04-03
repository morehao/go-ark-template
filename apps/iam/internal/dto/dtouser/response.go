package dtouser

import (
	"github.com/morehao/goark/apps/iam/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateResp struct {
	// UserID 用户UserID
	UserID uint `json:"userID"`
	// PersonID 自然人ID
	PersonID uint `json:"personID"`
}

type UserDetailResp struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objuser.UserBaseInfo
	gobject.OperatorBaseInfo
}

type UserPageListItem struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required"`
	objuser.UserBaseInfo
	gobject.OperatorBaseInfo
}

type UserPageListResp struct {
	// List 数据列表
	List []UserPageListItem `json:"list"`
	// Total 数据总条数
	Total int64 `json:"total"`
}

type UserDepartmentItem struct {
	DepartmentID   uint   `json:"departmentID"`
	DepartmentName string `json:"departmentName"`
	DeptType       string `json:"deptType"`
}

type UserDepartmentsResp struct {
	List []UserDepartmentItem `json:"list"`
}

// UserRoleItem 用户角色列表项
type UserRoleItem struct {
	// RoleID 角色ID
	RoleID uint `json:"roleId"`
	// RoleName 角色名称
	RoleName string `json:"roleName"`
	// RoleCode 角色编码
	RoleCode string `json:"roleCode"`
	// RoleType 角色类型
	RoleType string `json:"roleType"`
}

// UserRolesResp 用户角色列表响应
type UserRolesResp struct {
	// List 角色列表
	List []UserRoleItem `json:"list"`
}
