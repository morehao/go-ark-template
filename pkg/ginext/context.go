package ginext

import "github.com/gin-gonic/gin"

const (
	CtxKeyTenantID = "tenantId"
	CtxKeyOrgID    = "orgId"
	CtxKeyUserType = "userType"
	CtxKeyPersonID = "personId"
)

func GetTenantID(c *gin.Context) uint {
	return c.GetUint(CtxKeyTenantID)
}

func GetOrgID(c *gin.Context) uint {
	return c.GetUint(CtxKeyOrgID)
}

func GetUserType(c *gin.Context) string {
	return c.GetString(CtxKeyUserType)
}

func GetPersonID(c *gin.Context) uint {
	return c.GetUint(CtxKeyPersonID)
}
