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

// RegisterReq 用户注册请求
type RegisterReq struct {
	// TenantName 租户名称
	TenantName string `json:"tenantName" validate:"required" label:"租户名称"`
	// TenantCode 租户编码
	TenantCode string `json:"tenantCode" validate:"required" label:"租户编码"`
	// Username 用户名
	Username string `json:"username" validate:"required" label:"用户名"`
	// Password 密码
	Password string `json:"password" validate:"required" label:"密码"`
	// Mobile 手机号
	Mobile string `json:"mobile" label:"手机号"`
	// Email 邮箱
	Email string `json:"email" validate:"required" label:"邮箱"`
	// RealName 真实姓名
	RealName string `json:"realName" validate:"required" label:"真实姓名"`
}
