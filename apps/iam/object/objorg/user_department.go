package objorg

import "github.com/morehao/goark/apps/iam/iammodel"

type UserDepartmentBaseInfo struct {
	// TenantID 租户ID(冗余)
	TenantID uint `json:"tenantID" form:"tenantID"`
	// DeptID 部门ID
	DeptID uint `json:"deptID" form:"deptID"`
	// DeptType 部门类型: primary-主部门 secondary-其他部门
	DeptType iammodel.UserDeptType `json:"deptType" form:"deptType"`
	// UserID 用户ID
	UserID uint `json:"userID" form:"userID"`
}
