package objpermission

type RoleMenuBaseInfo struct {
	TenantID uint `json:"tenantID" form:"tenantID"` // 租户ID(冗余)
	MenuID   uint `json:"menuID" form:"menuID"`     // 菜单ID
	RoleID   uint `json:"roleID" form:"roleID"`     // 角色ID
}
