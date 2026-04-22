package objorg

type TenantApplicationBaseInfo struct {
	AppID    uint `json:"appID" form:"appID"`       // 应用ID
	TenantID uint `json:"tenantID" form:"tenantID"` // 租户ID
}
