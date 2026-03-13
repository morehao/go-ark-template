package code

import "github.com/morehao/golib/gerror"

const (
	TenantCreateError      = 100200
	TenantDeleteError      = 100201
	TenantUpdateError      = 100202
	TenantGetDetailError   = 100203
	TenantGetPageListError = 100204
	TenantNotExistError    = 100205
)

var tenantErrorMsgMap = gerror.CodeMsgMap{
	TenantCreateError:      "创建租户管理失败",
	TenantDeleteError:      "删除租户管理失败",
	TenantUpdateError:      "修改租户管理失败",
	TenantGetDetailError:   "查看租户管理失败",
	TenantGetPageListError: "查看租户管理列表失败",
	TenantNotExistError:    "租户管理不存在",
}
