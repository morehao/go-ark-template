package dtouser

import (
	"github.com/morehao/goark/apps/iam/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateResp struct {
	UserID   uint `json:"userID"`   // 用户UserID
	PersonID uint `json:"personID"` // 自然人ID
}

type UserDetailResp struct {
	ID uint `json:"id" validate:"required"` // 数据自增 ID
	objuser.UserBaseInfo
	gobject.OperatorBaseInfo
}

type UserPageListItem struct {
	ID uint `json:"id" validate:"required"` // 数据自增 ID
	objuser.UserBaseInfo
	gobject.OperatorBaseInfo
}

type UserPageListResp struct {
	List  []UserPageListItem `json:"list"`  // 数据列表
	Total int64              `json:"total"` // 数据总条数
}

type UserDepartmentItem struct {
	DepartmentID   uint   `json:"departmentID"`   // 部门ID
	DepartmentName string `json:"departmentName"` // 部门名称
	DeptType       string `json:"deptType"`       // 部门类型
}

type UserDepartmentsResp struct {
	List []UserDepartmentItem `json:"list"` // 用户部门列表
}

type UserRoleItem struct {
	RoleID   uint   `json:"roleId"`   // 角色ID
	RoleName string `json:"roleName"` // 角色名称
	RoleCode string `json:"roleCode"` // 角色编码
	RoleType string `json:"roleType"` // 角色类型
}

type UserRolesResp struct {
	List []UserRoleItem `json:"list"` // 角色列表
}
