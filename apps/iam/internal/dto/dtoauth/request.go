package dtoauth

type LoginReq struct {
	Username string `json:"username" validate:"required" label:"用户名/手机号/邮箱"`
	Password string `json:"password" validate:"required" label:"密码"`
}

type SelectTenantReq struct {
	PersonID uint `json:"personId" validate:"required" label:"自然人ID"`
	TenantID uint `json:"tenantId" validate:"required" label:"租户ID"`
}

type SwitchTenantReq struct {
	TenantID uint `json:"tenantId" validate:"required" label:"租户ID"`
}

type LogoutReq struct{}
