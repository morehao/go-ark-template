package dtoauth

type LoginResp struct {
	NeedSelectTenant bool           `json:"needSelectTenant"`
	Token            string         `json:"token,omitempty"`
	PersonID         uint           `json:"personId"`
	Tenants          []TenantItem   `json:"tenants,omitempty"`
	UserInfo         *LoginUserInfo `json:"userInfo,omitempty"`
}

type TenantItem struct {
	TenantID   uint   `json:"tenantId"`
	TenantName string `json:"tenantName"`
	TenantCode string `json:"tenantCode"`
	RoleName   string `json:"roleName,omitempty"`
}

type LoginUserInfo struct {
	UserID   uint   `json:"userId"`
	PersonID uint   `json:"personId"`
	TenantID uint   `json:"tenantId"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	UserType string `json:"userType"`
}

type SelectTenantResp struct {
	Token    string         `json:"token"`
	UserInfo *LoginUserInfo `json:"userInfo"`
}

type SwitchTenantResp = SelectTenantResp

type CurrentUserResp struct {
	UserID   uint   `json:"userId"`
	PersonID uint   `json:"personId"`
	TenantID uint   `json:"tenantId"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	UserType string `json:"userType"`
}
