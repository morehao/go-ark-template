package objorganization

type OrganizationBaseInfo struct {
	// Description 产品线描述
	Description string `json:"description" form:"description"`
	// SortOrder 排序
	SortOrder int32 `json:"sortOrder" form:"sortOrder"`
	// Status 状态: active-正常 inactive-停用
	Status string `json:"status" form:"status"`
	// OrganizationCode 产品线编码
	OrganizationCode string `json:"organizationCode" form:"organizationCode"`
	// OrganizationName 产品线名称
	OrganizationName string `json:"organizationName" form:"organizationName"`
}
