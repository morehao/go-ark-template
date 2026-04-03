package code

import "github.com/morehao/golib/gerror"

const (
	LoginFailedError          = 100800
	LogoutFailedError         = 100801
	PasswordIncorrectError    = 100802
	UserDisabledError         = 100803
	UserLockedError           = 100804
	TenantSelectRequiredError = 100805
	PersonNotFoundError       = 100806
	TokenBlacklistedError     = 100807
	RoleAssignMenusError      = 100808
	RoleRemoveMenusError      = 100809
	RoleListMenusError        = 100810
	UserAssignRolesError      = 100811
	UserRemoveRolesError      = 100812
	UserListRolesError        = 100813
	UserGetPermissionsError   = 100814
)

var authErrorMsgMap = gerror.CodeMsgMap{
	LoginFailedError:          "登录失败",
	LogoutFailedError:         "登出失败",
	PasswordIncorrectError:    "密码不正确",
	UserDisabledError:         "用户已被禁用",
	UserLockedError:           "用户已被锁定",
	TenantSelectRequiredError: "请选择租户",
	PersonNotFoundError:       "未找到对应用户",
	TokenBlacklistedError:     "令牌已失效",
	RoleAssignMenusError:      "分配菜单失败",
	RoleRemoveMenusError:      "移除菜单失败",
	RoleListMenusError:        "获取角色菜单列表失败",
	UserAssignRolesError:      "分配角色失败",
	UserRemoveRolesError:      "移除角色失败",
	UserListRolesError:        "获取用户角色列表失败",
	UserGetPermissionsError:   "获取用户权限失败",
}
