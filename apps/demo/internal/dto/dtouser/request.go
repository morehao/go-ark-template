package dtouser

import (
	"github.com/morehao/goark/demo/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateReq struct {
	objuser.UserBaseInfo
}

type UserUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
	objuser.UserBaseInfo
}

type UserDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
}

type UserPageListReq struct {
	gobject.PageQuery
}

type UserDeleteReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增 ID
}
