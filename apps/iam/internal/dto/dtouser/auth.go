package dtouser

import "github.com/morehao/goark/iam/model"

type LoginByPasswordReq struct {
	Account  string `json:"account" validate:"required" label:"登录账号"`
	Password string `json:"password" validate:"required" label:"登录密码"`
}

type SelectTenantReq struct {
	TenantID uint `json:"tenantId" validate:"required" label:"租户ID"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" validate:"required" label:"刷新令牌"`
}

type LogoutReq struct {
	RefreshToken string `json:"refreshToken" label:"刷新令牌"`
}

type RegisterReq struct {
	Username   string `json:"username" validate:"required" label:"用户名"`
	Password   string `json:"password" validate:"required" label:"密码"`
	Mobile     string `json:"mobile" label:"手机号"`
	Email      string `json:"email" validate:"required" label:"邮箱"`
	RealName   string `json:"realName" validate:"required" label:"真实姓名"`
	InviteCode string `json:"inviteCode" label:"邀请码"`
	SSOType    string `json:"ssoType" label:"SSO类型"`
	OpenID     string `json:"openID" label:"OpenID"`
}

type ApproveReq struct {
	UserID   uint   `json:"userId" validate:"required" label:"用户ID"`
	Approved bool   `json:"approved" label:"是否通过"`
}

type UnlockAccountReq struct {
	Account   string `json:"account" validate:"required" label:"账号"`
	Captcha   string `json:"captcha" validate:"required" label:"验证码"`
	CaptchaID string `json:"captchaId" label:"验证码ID"`
}

type LoginByPasswordResp struct {
	TempToken        string           `json:"tempToken"`
	NeedSelectTenant bool             `json:"needSelectTenant"`
	TenantList       []TenantListItem `json:"tenantList"`
	PersonID         uint             `json:"personId"`
	RealName         string           `json:"realName"`
}

type TenantListItem struct {
	TenantID   uint   `json:"tenantId"`
	TenantName string `json:"tenantName"`
	OrgID      uint   `json:"orgId"`
	OrgName    string `json:"orgName"`
}

type SelectTenantResp struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
	UserInfo     LoginUserInfo `json:"userInfo"`
}

type RefreshTokenResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUserInfo struct {
	UserID     uint           `json:"userId"`
	PersonID   uint           `json:"personId"`
	Username   string         `json:"username"`
	RealName   string         `json:"realName"`
	UserType   model.UserType `json:"userType"`
	TenantID   uint           `json:"tenantId"`
	TenantName string         `json:"tenantName"`
}

type RegisterResp struct {
	UserID       uint   `json:"userId"`
	PersonID     uint   `json:"personId"`
	Status       string `json:"status"`
	PersonExists bool   `json:"personExists"`
	Message      string `json:"message"`
}