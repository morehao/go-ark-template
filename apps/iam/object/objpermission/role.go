package objpermission

type RoleBaseInfo struct {
	// TenantID 所属租户ID
	TenantID uint `json:"tenantID" form:"tenantID"`
	// DataScope 数据权限范围: all-全部 dept_and_sub-本部门及以下 dept-本部门 self-仅本人 custom-自定义
	DataScope string `json:"dataScope" form:"dataScope"`
	// Description 角色描述
	Description string `json:"description" form:"description"`
	// RoleCode 角色编码
	RoleCode string `json:"roleCode" form:"roleCode"`
	// RoleName 角色名称
	RoleName string `json:"roleName" form:"roleName"`
	// RoleType 角色类型: custom-自定义 system-系统内置
	RoleType string `json:"roleType" form:"roleType"`
	// SortOrder 排序
	SortOrder int32 `json:"sortOrder" form:"sortOrder"`
	// Status 状态: active-正常 inactive-停用
	Status string `json:"status" form:"status"`
}
