package objorg

type OrganizationBaseInfo struct {
	Domain           string `json:"domain" form:"domain"`
	Logo             string `json:"logo" form:"logo"`
	Description      string `json:"description" form:"description"`
	SortOrder        int32  `json:"sortOrder" form:"sortOrder"`
	Status           string `json:"status" form:"status"`
	OrganizationCode string `json:"organizationCode" form:"organizationCode"`
	OrganizationName string `json:"organizationName" form:"organizationName"`
}

type OrganizationAdminInfo struct {
	Username string `json:"username" form:"username" validate:"required" label:"管理员用户名"`
	Mobile   string `json:"mobile" form:"mobile" label:"管理员手机号"`
	Email    string `json:"email" form:"email" label:"管理员邮箱"`
	RealName string `json:"realName" form:"realName" validate:"required" label:"管理员真实姓名"`
	Password string `json:"password" form:"password" label:"管理员密码"`
}
