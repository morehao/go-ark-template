package objorg

type TenantBaseInfo struct {
	Address                 string `json:"address" form:"address"`                                 // 地址
	ContactEmail            string `json:"contactEmail" form:"contactEmail"`                       // 联系邮箱
	ContactPhone            string `json:"contactPhone" form:"contactPhone"`                       // 联系电话
	LegalPerson             string `json:"legalPerson" form:"legalPerson"`                         // 法人
	Logo                    string `json:"logo" form:"logo"`                                       // Logo
	OrganizationID          uint   `json:"organizationID" form:"organizationID"`                   // 组织ID
	ShortName               string `json:"shortName" form:"shortName"`                             // 简称
	Status                  string `json:"status" form:"status"`                                   // 状态
	TenantCode              string `json:"tenantCode" form:"tenantCode"`                           // 租户编码
	TenantName              string `json:"tenantName" form:"tenantName"`                           // 租户名称
	UnifiedSocialCreditCode string `json:"unifiedSocialCreditCode" form:"unifiedSocialCreditCode"` // 统一社会信用代码
}

type TenantAdminInfo struct {
	Username string `json:"username" form:"username" validate:"required" label:"管理员用户名"`  // 管理员用户名
	Mobile   string `json:"mobile" form:"mobile" label:"管理员手机号"`                          // 管理员手机号
	Email    string `json:"email" form:"email" label:"管理员邮箱"`                             // 管理员邮箱
	RealName string `json:"realName" form:"realName" validate:"required" label:"管理员真实姓名"` // 管理员真实姓名
}
