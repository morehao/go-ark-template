package dtouser

import (
	"github.com/morehao/goark/apps/demo/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateResp struct {
	ID uint `json:"id"` // 数据自增 ID
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
