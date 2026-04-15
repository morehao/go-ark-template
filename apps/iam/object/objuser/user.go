package objuser

type UserBaseInfo struct {
	TenantID    uint   `json:"tenantID" form:"tenantID"`       // 所属租户ID, 0表示平台管理员账号
	DeptID      uint   `json:"deptID" form:"deptID"`           // 主部门ID(冗余字段,方便查询,实际关联关系在iam_user_department表)
	EmployeeNo  string `json:"employeeNo" form:"employeeNo"`   // 工号
	EntryDate   int64  `json:"entryDate" form:"entryDate"`     // 入职日期
	JobLevel    string `json:"jobLevel" form:"jobLevel"`       // 职级
	LastLoginAt int64  `json:"lastLoginAt" form:"lastLoginAt"` // 最后登录时间
	LastLoginIp string `json:"lastLoginIp" form:"lastLoginIp"` // 最后登录IP(支持IPv6)
	LoginCount  int32  `json:"loginCount" form:"loginCount"`   // 登录次数
	PersonID    uint   `json:"personID" form:"personID"`       // 自然人ID
	Position    string `json:"position" form:"position"`       // 职位
	Status      string `json:"status" form:"status"`           // 状态: enabled-正常 locked-锁定 disabled-禁用
	UserType    string `json:"userType" form:"userType"`       // 用户类型: normal-普通用户 tenant_admin-租户管理员 platform_admin-平台管理员(可管理所有租户)
	Username    string `json:"username" form:"username"`       // 用户名(租户用户:租户内唯一,平台管理员:全局唯一,需应用层保证)
}
