package dtouser

import (
	"github.com/morehao/goark/apps/demo/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateReq struct {
	objuser.UserBaseInfo
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
