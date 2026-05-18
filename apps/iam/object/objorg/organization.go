package objorg

type OrgBaseInfo struct {
	DisplayCode string `json:"displayCode" form:"displayCode"`           // 组织编码(对外展示)
	OrgName     string `json:"orgName" form:"orgName" validate:"required"` // 组织名称
	Description string `json:"description" form:"description"`           // 描述
	Logo        string `json:"logo" form:"logo"`                         // Logo
	Sequence    int32  `json:"sequence" form:"sequence"`               // 排序
	Status      string `json:"status" form:"status"`                     // 状态
}

type OrgAdminInfo struct {
	Username string `json:"username" form:"username" validate:"required" label:"管理员用户名"`  // 管理员用户名
	Mobile   string `json:"mobile" form:"mobile" validate:"required" label:"管理员手机号"`     // 管理员手机号
	Email    string `json:"email" form:"email" validate:"required" label:"管理员邮箱"`        // 管理员邮箱
	RealName string `json:"realName" form:"realName" validate:"required" label:"管理员真实姓名"` // 管理员真实姓名
	Password string `json:"password" form:"password" label:"管理员密码"`                       // 管理员密码
}