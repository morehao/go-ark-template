package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/morehao/goark/apps/iam/docs"

	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/apps/iam/internal/router"
	"github.com/morehao/goark/apps/iam/internal/service/svcauth"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gcontext"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gmiddleware/ginmiddleware"
	"github.com/morehao/golib/biz/grouter/ginrouter"
	"github.com/morehao/golib/gerror"
)

const AppName = "iam"

func Routers(engine *gin.Engine) {
	engine.Use(ginmiddleware.AccessLog())
	engine.Use(ginmiddleware.JWTAuth(config.Conf.JWT.SignKey, ginmiddleware.WithWhiteList([]string{
		"/iam/redocs",
		"/iam/docs/*any",
		"/iam/docs/doc.json",
		"/v1/iam/auth/login",
	})))
	if config.Conf.Server.Env == "dev" {
		swaggerGroup := engine.Group("/" + AppName)
		ginrouter.RegisterSwagger(swaggerGroup, AppName)
	}

	v1Group := engine.Group(fmt.Sprintf("%s/%s", gconstant.ApiVersionV1, AppName))

	routerGroups := &ginrouter.RouterGroups{
		V1: v1Group,
	}
	router.RegisterRouter(routerGroups, AppName)
}

// tokenBlacklistCheck token黑名单检查中间件
func tokenBlacklistCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		if svcauth.IsTokenBlacklisted(c, token) {
			gincontext.Abort(c, &gerror.Error{
				Code: gconstant.TokenInvalidErr,
				Msg:  "token已失效",
			})
			return
		}
		c.Next()
	}
}

// ginContextToStdContext 将gin context中的用户信息传播到标准context
func ginContextToStdContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request != nil {
			ctx := c.Request.Context()
			if tenantID := gincontext.GetTenantID(c); tenantID > 0 {
				ctx = context.WithValue(ctx, gcontext.KeyTenantID, tenantID)
			}
			if orgID := gincontext.GetOrgID(c); orgID > 0 {
				ctx = context.WithValue(ctx, gcontext.KeyOrgID, orgID)
			}
			if userType := gincontext.GetUserType(c); userType != "" {
				ctx = context.WithValue(ctx, gcontext.KeyUserType, userType)
			}
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
