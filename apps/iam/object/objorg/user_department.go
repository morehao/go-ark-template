package objorg

type UserDepartmentBaseInfo struct {
	TenantID uint   `json:"tenantID" form:"tenantID"` // 租户ID(冗余)
	DeptID   uint   `json:"deptID" form:"deptID"`     // 部门ID
	DeptType string `json:"deptType" form:"deptType"` // 部门类型: primary-主部门 secondary-其他部门
	UserID   uint   `json:"userID" form:"userID"`     // 用户ID
}
