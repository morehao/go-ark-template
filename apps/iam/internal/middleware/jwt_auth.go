package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/pkg/contextkeys"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/gauth/jwtauth"
)

type JWTCustomData struct {
	UserID   uint   `json:"userId"`
	PersonID uint   `json:"personId"`
	TenantID uint   `json:"tenantId"`
	OrgID    uint   `json:"orgId"`
	UserType string `json:"userType"`
}

var noAuthPaths = []string{
	"/auth/login",
	"/auth/selectTenant",
	"/swagger/",
}

func shouldSkipAuth(path string) bool {
	for _, p := range noAuthPaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

func JWTAuth(signKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if shouldSkipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			gincontext.Abort(c, code.GetError(gconstant.TokenInvalidErr))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			gincontext.Abort(c, code.GetError(gconstant.TokenInvalidErr))
			return
		}

		if IsBlacklisted(c.Request.Context(), tokenStr) {
			gincontext.Abort(c, code.GetError(gconstant.TokenExpiredErr))
			return
		}

		var claims jwtauth.Claims[JWTCustomData]
		if err := jwtauth.ParseToken(signKey, tokenStr, &claims); err != nil {
			gincontext.Abort(c, code.GetError(gconstant.TokenInvalidErr))
			return
		}

		c.Set(string(contextkeys.KeyUserID), claims.CustomData.UserID)
		c.Set(string(contextkeys.KeyTenantID), claims.CustomData.TenantID)
		c.Set(string(contextkeys.KeyOrgID), claims.CustomData.OrgID)
		c.Set(string(contextkeys.KeyUserType), claims.CustomData.UserType)
		c.Set(string(contextkeys.KeyPersonID), claims.CustomData.PersonID)

		c.Next()
	}
}
