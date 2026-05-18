package dtoapikey

import "github.com/morehao/golib/biz/gobject"

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
	ID uint `json:"id" form:"id" binding:"required"` // ID
}

type ApiKeyListReq struct {
	gobject.PageQuery
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"`  // 每页数据条数
	AppID    uint   `json:"appID" form:"appID"`       // 应用ID
	KeyName  string `json:"keyName" form:"keyName"`   // 密钥名称
	Status   string `json:"status" form:"status"`     // 状态
}

type ApiKeyListResp struct {
	List  []ApiKeyListItem `json:"list"`  // 数据列表
	Total int64            `json:"total"` // 数据总条数
}

type ApiKeyListItem struct {
	ID          uint   `json:"id"`           // ID
	AppID       uint   `json:"appID"`        // 应用ID
	KeyName     string `json:"keyName"`      // 密钥名称
	KeyPrefix   string `json:"keyPrefix"`    // 密钥前缀
	ApiKey      string `json:"apiKey"`       // 解密后的完整 API Key
	Scopes      string `json:"scopes"`       // 权限范围
	AccessPolicy string `json:"accessPolicy"` // 访问策略
	AllowedIPs  string `json:"allowedIPs"`  // 允许的IP列表
	Status      string `json:"status"`       // 状态
	LastUsedAt  string `json:"lastUsedAt"`   // 最后使用时间
	ExpiresAt   string `json:"expiresAt"`    // 过期时间
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
}

type ApiKeyDisableReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // ID
}

type ApiKeyEnableReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // ID
}
