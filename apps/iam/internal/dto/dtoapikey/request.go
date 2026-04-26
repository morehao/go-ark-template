package dtoapikey

type ApiKeyCreateReq struct {
	AppID       uint   `json:"appID" comment:"应用ID"`
	KeyName     string `json:"keyName" comment:"密钥名称" binding:"required"`
	Scopes      string `json:"scopes" comment:"权限范围: user:read,user:write"`
	AccessPolicy string `json:"accessPolicy" comment:"访问策略: all-全部 ip-限IP"`
	AllowedIPs  string `json:"allowedIPs" comment:"允许的IP列表(JSON数组)"`
	ExpiresAt   string `json:"expiresAt" comment:"过期时间(YYYY-MM-DD HH:mm:ss)，空表示永不过期"`
}

type ApiKeyCreateResp struct {
	ID                 uint   `json:"id"`
	KeyPrefix          string `json:"keyPrefix" comment:"密钥前缀"`
	ApiKey             string `json:"apiKey" comment:"完整 API Key 明文"`
	ExpiresAt          string `json:"expiresAt" comment:"过期时间"`
}

type ApiKeyDeleteReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type ApiKeyListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	AppID    uint   `json:"appID" form:"appID"`
	KeyName  string `json:"keyName" form:"keyName"`
	Status   string `json:"status" form:"status"`
}

type ApiKeyListResp struct {
	List  []ApiKeyListItem `json:"list"`
	Total int64           `json:"total"`
}

type ApiKeyListItem struct {
	ID                 uint   `json:"id"`
	AppID              uint   `json:"appID"`
	KeyName            string `json:"keyName"`
	KeyPrefix          string `json:"keyPrefix"`
	ApiKey             string `json:"apiKey" comment:"解密后的完整 API Key"`
	Scopes             string `json:"scopes"`
	AccessPolicy       string `json:"accessPolicy"`
	AllowedIPs         string `json:"allowedIPs"`
	Status             string `json:"status"`
	LastUsedAt         string `json:"lastUsedAt"`
	ExpiresAt          string `json:"expiresAt"`
	CreatedAt          int64  `json:"createdAt"`
}

type ApiKeyDisableReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type ApiKeyEnableReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}
