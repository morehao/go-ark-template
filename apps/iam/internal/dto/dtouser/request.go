package dtouser

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/apps/iam/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateReq struct {
	objuser.UserBaseInfo
	// Mobile 手机号
	Mobile string `json:"mobile" form:"mobile"`
	// Email 邮箱
	Email string `json:"email" form:"email"`
	// RealName 真实姓名
	RealName string `json:"realName" form:"realName"`
	// PrimaryDeptID 主部门ID，未传入时自动使用租户顶级部门
	PrimaryDeptID uint `json:"primaryDeptID" form:"primaryDeptID"`
	// SecondaryDeptIDs 其他关联部门ID列表（第二部门）
	SecondaryDeptIDs []uint `json:"secondaryDeptIDs" form:"secondaryDeptIDs"`
}

type UserUpdateReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	objuser.UserBaseInfo
}

type UserDetailReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type UserPageListReq struct {
	gobject.PageQuery
}

type UserDeleteReq struct {
	// ID 数据自增 ID
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type UserDepartmentAssignReq struct {
	UserID       uint                  `json:"userID" validate:"required" label:"用户ID"`
	DepartmentID uint                  `json:"departmentID" validate:"required" label:"部门ID"`
	DeptType     iammodel.UserDeptType `json:"deptType" validate:"required" label:"部门类型: primary-secondary"`
}

type UserDepartmentRemoveReq struct {
	UserID       uint `json:"userID" validate:"required" label:"用户ID"`
	DepartmentID uint `json:"departmentID" validate:"required" label:"部门ID"`
}

type UserDepartmentsReq struct {
	UserID uint `json:"userID" form:"userID" validate:"required" label:"用户ID"`
}

// UserAssignRolesReq 用户分配角色请求
type UserAssignRolesReq struct {
	// UserID 用户ID
	UserID uint `json:"userId" validate:"required" label:"用户ID"`
	// RoleIDs 角色ID列表
	RoleIDs []uint `json:"roleIds" validate:"required" label:"角色ID列表"`
}

// UserRemoveRolesReq 用户移除角色请求
type UserRemoveRolesReq struct {
	// UserID 用户ID
	UserID uint `json:"userId" validate:"required" label:"用户ID"`
	// RoleIDs 角色ID列表
	RoleIDs []uint `json:"roleIds" validate:"required" label:"角色ID列表"`
}

// UserRolesReq 查询用户角色请求
type UserRolesReq struct {
	// UserID 用户ID
	UserID uint `json:"userId" form:"userId" validate:"required" label:"用户ID"`
}
