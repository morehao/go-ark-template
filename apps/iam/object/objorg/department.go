package objorg

type DepartmentBaseInfo struct {
	TenantID  uint   `json:"tenantID" form:"tenantID"`   // 所属租户ID
	DeptCode  string `json:"deptCode" form:"deptCode"`   // 部门编码
	DeptLevel int32  `json:"deptLevel" form:"deptLevel"` // 部门层级
	DeptName  string `json:"deptName" form:"deptName"`   // 部门名称
	DeptPath  string `json:"deptPath" form:"deptPath"`   // 部门路径: /1/2/3/
	LeaderID  uint   `json:"leaderID" form:"leaderID"`   // 部门负责人ID
	ParentID  uint   `json:"parentID" form:"parentID"`   // 父部门ID,0表示根部门
	Sequence int32  `json:"sequence" form:"sequence"` // 排序
	Status    string `json:"status" form:"status"`       // 状态: enabled-启用 disabled-停用
}
