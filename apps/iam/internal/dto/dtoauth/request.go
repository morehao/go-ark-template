package dtoauth

// LoginByPasswordReq 密码登录请求
type LoginByPasswordReq struct {
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

// RefreshTokenReq 刷新令牌请求
type RefreshTokenReq struct {
	// RefreshToken 刷新令牌
	RefreshToken string `json:"refreshToken" validate:"required" label:"刷新令牌"`
}

// LogoutReq 登出请求
type LogoutReq struct {
	// RefreshToken 刷新令牌（可选，传递后一并吊销）
	RefreshToken string `json:"refreshToken" label:"刷新令牌"`
}
