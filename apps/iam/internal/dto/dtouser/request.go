package dtouser

import (
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
