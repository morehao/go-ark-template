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
