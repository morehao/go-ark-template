package dtoauth

// LoginResp 登录响应
type LoginResp struct {
	// Token JWT令牌(单租户时直接返回完整token)
	Token string `json:"token"`
	// NeedSelectTenant 是否需要选择租户
	NeedSelectTenant bool `json:"needSelectTenant"`
	// Tenants 可选租户列表(多租户时返回)
	Tenants []TenantItem `json:"tenants"`
	// UserInfo 用户信息(单租户时返回)
	UserInfo *LoginUserInfo `json:"userInfo"`
}

// TenantItem 租户列表项
type TenantItem struct {
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
	// UserInfo 用户信息
	UserInfo *LoginUserInfo `json:"userInfo"`
}

// LoginUserInfo 登录用户信息
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
	UserType string `json:"userType"`
	// TenantID 租户ID
	TenantID uint `json:"tenantId"`
	// TenantName 租户名称
	TenantName string `json:"tenantName"`
}
