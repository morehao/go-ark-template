package objorg

type OrganizationApplicationBaseInfo struct {
	AppID uint `json:"appID" form:"appID"` // 应用ID
	OrgID uint `json:"orgID" form:"orgID"` // 组织ID
}
