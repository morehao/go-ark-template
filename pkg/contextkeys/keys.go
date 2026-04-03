package contextkeys

type contextKey string

const (
	KeyUserID   contextKey = "userId"
	KeyTenantID contextKey = "tenantId"
	KeyOrgID    contextKey = "orgId"
	KeyDeptID   contextKey = "deptId"
	KeyPersonID contextKey = "personId"
	KeyUserType contextKey = "userType"
)
