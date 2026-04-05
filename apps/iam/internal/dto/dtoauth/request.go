package dtoauth

// LoginReq 登录请求
type LoginReq struct {
	// Account 登录账号(手机号或邮箱)
	Account string `json:"account" validate:"required" label:"登录账号"`
	// Password 登录密码
	Password string `json:"password" validate:"required" label:"登录密码"`
}

// SelectTenantReq 选择租户请求
type SelectTenantReq struct {
	// TenantID 租户ID
	TenantID uint `json:"tenantId" validate:"required" label:"租户ID"`
}
