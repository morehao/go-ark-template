package dtoauth

import "github.com/morehao/goark/apps/iam/model"

type LoginByPasswordResp struct {
	TempToken        string           `json:"tempToken"`        // JWT临时令牌(需通过selectTenant换取正式token)
	NeedSelectTenant bool             `json:"needSelectTenant"` // 是否需要选择租户
	TenantList       []TenantListItem `json:"tenantList"`       // 可选租户列表
	PersonID         uint             `json:"personId"`         // 自然人ID
	RealName         string           `json:"realName"`         // 真实姓名
}

type TenantListItem struct {
	TenantID   uint   `json:"tenantId"`   // 租户ID
	TenantName string `json:"tenantName"` // 租户名称
	OrgID      uint   `json:"orgId"`      // 组织ID
	OrgName    string `json:"orgName"`    // 组织名称
}

type SelectTenantResp struct {
	Token        string        `json:"token"`        // JWT令牌
	RefreshToken string        `json:"refreshToken"` // 刷新令牌
	UserInfo     LoginUserInfo `json:"userInfo"`     // 用户信息
}

type RefreshTokenResp struct {
	Token        string `json:"token"`        // 新的JWT令牌
	RefreshToken string `json:"refreshToken"` // 新的刷新令牌
}

type LoginUserInfo struct {
	UserID     uint              `json:"userId"`     // 用户ID
	PersonID   uint              `json:"personId"`   // 自然人ID
	Username   string            `json:"username"`   // 用户名
	RealName   string            `json:"realName"`   // 真实姓名
	UserType   model.UserType `json:"userType"`   // 用户类型
	TenantID   uint              `json:"tenantId"`   // 租户ID
	TenantName string            `json:"tenantName"` // 租户名称
}

type RegisterResp struct {
	UserID       uint   `json:"userId"`       // 用户ID
	PersonID     uint   `json:"personId"`     // 自然人ID
	Status       string `json:"status"`       // 状态
	PersonExists bool   `json:"personExists"` // Person是否已存在
	Message      string `json:"message"`      // 提示信息
}
