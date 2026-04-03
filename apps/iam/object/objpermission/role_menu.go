package objpermission

type RoleMenuBaseInfo struct {
	// TenantID 租户ID(冗余)
	TenantID uint `json:"tenantID" form:"tenantID"`
	// MenuID 菜单ID
	MenuID uint `json:"menuID" form:"menuID"`
	// RoleID 角色ID
	RoleID uint `json:"roleID" form:"roleID"`
}
