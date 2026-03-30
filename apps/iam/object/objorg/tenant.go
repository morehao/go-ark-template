package objorg

type TenantBaseInfo struct {
	Address                 string `json:"address" form:"address"`
	ContactEmail            string `json:"contactEmail" form:"contactEmail"`
	ContactPhone            string `json:"contactPhone" form:"contactPhone"`
	LegalPerson             string `json:"legalPerson" form:"legalPerson"`
	Logo                    string `json:"logo" form:"logo"`
	OrganizationID          uint   `json:"organizationID" form:"organizationID"`
	ShortName               string `json:"shortName" form:"shortName"`
	Status                  string `json:"status" form:"status"`
	TenantCode              string `json:"tenantCode" form:"tenantCode"`
	TenantName              string `json:"tenantName" form:"tenantName"`
	UnifiedSocialCreditCode string `json:"unifiedSocialCreditCode" form:"unifiedSocialCreditCode"`
}

type TenantAdminInfo struct {
	Username string `json:"username" form:"username" validate:"required" label:"管理员用户名"`
	Mobile   string `json:"mobile" form:"mobile" label:"管理员手机号"`
	Email    string `json:"email" form:"email" label:"管理员邮箱"`
	RealName string `json:"realName" form:"realName" validate:"required" label:"管理员真实姓名"`
}
