package objpermission

type UserRoleBaseInfo struct {
	TenantID uint `json:"tenantID" form:"tenantID"` // 租户ID(冗余)
	RoleID   uint `json:"roleID" form:"roleID"`     // 角色ID
	UserID   uint `json:"userID" form:"userID"`     // 用户ID
}
