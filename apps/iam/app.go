package iam

import (
	"context"

	"github.com/gin-gonic/gin"
	_ "github.com/morehao/goark/apps/iam/docs"
	"github.com/morehao/goark/apps/iam/internal/middleware"
	"github.com/morehao/goark/apps/iam/router"
	"github.com/morehao/goark/pkg/contextkeys"
	"github.com/morehao/goark/pkg/ginext"
)

const AppName = "iam"

// ginContextToStdContext transfers auth info from gin.Context to std context
// so that GORM organization callbacks can read tenant/org scope.
func ginContextToStdContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request != nil {
			ctx := c.Request.Context()
			if tenantID := ginext.GetTenantID(c); tenantID > 0 {
				ctx = context.WithValue(ctx, contextkeys.KeyTenantID, tenantID)
			}
			if orgID := ginext.GetOrgID(c); orgID > 0 {
				ctx = context.WithValue(ctx, contextkeys.KeyOrgID, orgID)
			}
			if userType := ginext.GetUserType(c); userType != "" {
				ctx = context.WithValue(ctx, contextkeys.KeyUserType, userType)
			}
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func Routers(engine *gin.Engine, signKey string) {
	routerGroup := engine.Group(AppName)

	noAuthGroup := routerGroup.Group("")

	authGroup := routerGroup.Group("")
	authGroup.Use(middleware.JWTAuth(signKey))
	authGroup.Use(ginContextToStdContext())

	routerGroups := &router.RouterGroups{
		AuthGroup:   authGroup,
		NoAuthGroup: noAuthGroup,
	}
	router.RegisterRouter(routerGroups, AppName)
}
