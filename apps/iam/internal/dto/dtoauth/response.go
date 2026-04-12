package dtoauth

import "github.com/morehao/goark/apps/iam/iammodel"

// LoginByPasswordResp 密码登录响应
type LoginByPasswordResp struct {
	// TempToken JWT临时令牌(需通过selectTenant换取正式token)
	TempToken string `json:"tempToken"`
	// NeedSelectTenant 是否需要选择租户
	NeedSelectTenant bool `json:"needSelectTenant"`
	// TenantList 可选租户列表
	TenantList []TenantListItem `json:"tenantList"`
	// PersonID 自然人ID
	PersonID uint `json:"personId"`
	// RealName 真实姓名
	RealName string `json:"realName"`
}

// TenantListItem 租户列表项
type TenantListItem struct {
	// TenantID 租户ID
	TenantID uint `json:"tenantId"`
	// TenantName 租户名称
	TenantName string `json:"tenantName"`
	// OrgID 产品线ID
	OrgID uint `json:"orgId"`
	// OrgName 产品线名称
	OrgName string `json:"orgName"`
}

// SelectTenantResp 选择租户响应
type SelectTenantResp struct {
	// Token JWT令牌
	Token string `json:"token"`
	// RefreshToken 刷新令牌
	RefreshToken string `json:"refreshToken"`
	// UserInfo 用户信息
	UserInfo LoginUserInfo `json:"userInfo"`
}

// RefreshTokenResp 刷新令牌响应
type RefreshTokenResp struct {
	// Token 新的JWT令牌
	Token string `json:"token"`
	// RefreshToken 新的刷新令牌
	RefreshToken string `json:"refreshToken"`
}

// LoginUserInfo 密码登录用户信息
type LoginUserInfo struct {
	// UserID 用户ID
	UserID uint `json:"userId"`
	// PersonID 自然人ID
	PersonID uint `json:"personId"`
	// Username 用户名
	Username string `json:"username"`
	// RealName 真实姓名
	RealName string `json:"realName"`
	// UserType 用户类型
	UserType iammodel.UserType `json:"userType"`
	// TenantID 租户ID
	TenantID uint `json:"tenantId"`
	// TenantName 租户名称
	TenantName string `json:"tenantName"`
}

// RegisterResp 用户注册响应
type RegisterResp struct {
	// TenantID 租户ID
	TenantID uint `json:"tenantId"`
	// UserID 用户ID
	UserID uint `json:"userId"`
	// PersonID 自然人ID
	PersonID uint `json:"personId"`
	// Status 用户状态
	Status string `json:"status"`
	// PersonExists Person是否已存在
	PersonExists bool `json:"personExists"`
	// Message 提示信息
	Message string `json:"message"`
}
