package code

import "github.com/morehao/golib/gerror"

const (
	AuthLoginError           = 100800
	AuthLogoutError          = 100801
	AuthPasswordError        = 100802
	AuthAccountDisabledError = 100803
	AuthTenantSelectError    = 100804
	AuthTokenBlacklistError  = 100805
	AuthPersonNotFoundError  = 100806
	AuthNoTenantError        = 100807
	AuthTokenGenerateError   = 100808
)

var authErrorMsgMap = gerror.CodeMsgMap{
	AuthLoginError:           "登录失败",
	AuthLogoutError:          "登出失败",
	AuthPasswordError:        "账号或密码错误",
	AuthAccountDisabledError: "账号已被禁用",
	AuthTenantSelectError:    "选择租户失败",
	AuthTokenBlacklistError:  "token加入黑名单失败",
	AuthPersonNotFoundError:  "用户不存在",
	AuthNoTenantError:        "该用户未关联任何租户",
	AuthTokenGenerateError:   "生成token失败",
}
