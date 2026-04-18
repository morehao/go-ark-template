package dtoauth

type LoginByPasswordReq struct {
	Account  string `json:"account" validate:"required" label:"登录账号"`  // 登录账号(手机号或邮箱)
	Password string `json:"password" validate:"required" label:"登录密码"` // 登录密码
}

type SelectTenantReq struct {
	TenantID uint `json:"tenantId" validate:"required" label:"租户ID"` // 租户ID
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" validate:"required" label:"刷新令牌"` // 刷新令牌
}

type LogoutReq struct {
	RefreshToken string `json:"refreshToken" label:"刷新令牌"` // 刷新令牌（可选，传递后一并吊销）
}

type RegisterReq struct {
	TenantName string `json:"tenantName" validate:"required" label:"租户名称"` // 租户名称
	TenantCode string `json:"tenantCode" validate:"required" label:"租户编码"` // 租户编码
	Username   string `json:"username" validate:"required" label:"用户名"`    // 用户名
	Password   string `json:"password" validate:"required" label:"密码"`     // 密码
	Mobile     string `json:"mobile" label:"手机号"`                          // 手机号
	Email      string `json:"email" validate:"required" label:"邮箱"`        // 邮箱
	RealName   string `json:"realName" validate:"required" label:"真实姓名"`   // 真实姓名
}
