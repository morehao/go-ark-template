package objpermission

type UserRoleBaseInfo struct {
	// TenantID 租户ID(冗余)
	TenantID uint `json:"tenantID" form:"tenantID"`
	// RoleID 角色ID
	RoleID uint `json:"roleID" form:"roleID"`
	// UserID 用户ID
	UserID uint `json:"userID" form:"userID"`
}
