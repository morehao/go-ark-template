package ginext

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/contextkeys"
)

func GetTenantID(c *gin.Context) uint {
	return c.GetUint(string(contextkeys.KeyTenantID))
}

func GetOrgID(c *gin.Context) uint {
	return c.GetUint(string(contextkeys.KeyOrgID))
}

func GetUserType(c *gin.Context) string {
	return c.GetString(string(contextkeys.KeyUserType))
}

func GetPersonID(c *gin.Context) uint {
	return c.GetUint(string(contextkeys.KeyPersonID))
}
