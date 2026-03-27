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
