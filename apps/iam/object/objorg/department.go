package objorg

type DepartmentBaseInfo struct {
	// TenantID 所属租户ID
	TenantID uint `json:"tenantID" form:"tenantID"`
	// DeptCode 部门编码
	DeptCode string `json:"deptCode" form:"deptCode"`
	// DeptLevel 部门层级
	DeptLevel int32 `json:"deptLevel" form:"deptLevel"`
	// DeptName 部门名称
	DeptName string `json:"deptName" form:"deptName"`
	// DeptPath 部门路径: /1/2/3/
	DeptPath string `json:"deptPath" form:"deptPath"`
	// LeaderID 部门负责人ID
	LeaderID uint `json:"leaderID" form:"leaderID"`
	// ParentID 父部门ID,0表示根部门
	ParentID uint `json:"parentID" form:"parentID"`
	// SortOrder 排序
	SortOrder int32 `json:"sortOrder" form:"sortOrder"`
	// Status 状态: enabled-启用 disabled-停用
	Status string `json:"status" form:"status"`
}
