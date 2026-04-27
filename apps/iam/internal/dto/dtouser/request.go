package dtouser

import (
	"github.com/morehao/goark/apps/iam/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateReq struct {
	objuser.UserBaseInfo
	Mobile           string `json:"mobile" form:"mobile"`                     // 手机号
	Email            string `json:"email" form:"email"`                       // 邮箱
	RealName         string `json:"realName" form:"realName"`                 // 真实姓名
	PrimaryDeptID    uint   `json:"primaryDeptID" form:"primaryDeptID"`       // 主部门ID，未传入时自动使用租户顶级部门
	SecondaryDeptIDs []uint `json:"secondaryDeptIDs" form:"secondaryDeptIDs"` // 其他关联部门ID列表（第二部门）
}

type UserUpdateReq struct {
	UserID uint `json:"userID" validate:"required" label:"数据自增id"` // 数据自增 ID
	objuser.UserBaseInfo
}

type UserDetailReq struct {
	UserID uint `json:"userID" form:"userID" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type UserPageListReq struct {
	gobject.PageQuery
}

type UserDeleteReq struct {
	UserID uint `json:"userID" form:"userID" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type UserDepartmentsAssignReq struct {
	UserID           uint   `json:"userID" validate:"required" label:"用户ID"`         // 用户ID
	PrimaryDeptID    uint   `json:"primaryDeptID" validate:"required" label:"主部门ID"` // 主部门ID
	SecondaryDeptIDs []uint `json:"secondaryDeptIDs" label:"其他部门ID列表"`               // 其他部门ID列表
}

type UserDepartmentsReq struct {
	UserID uint `json:"userID" form:"userID" validate:"required" label:"用户ID"` // 用户ID
}

type UserAssignRolesReq struct {
	UserID  uint   `json:"userId" validate:"required" label:"用户ID"`    // 用户ID
	RoleIDs []uint `json:"roleIds" validate:"required" label:"角色ID列表"` // 角色ID列表
}

type UserRolesReq struct {
	UserID uint `json:"userId" form:"userId" validate:"required" label:"用户ID"` // 用户ID
}

type UpdateProfileReq struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type LoginHistoryReq struct {
	gobject.PageQuery
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
