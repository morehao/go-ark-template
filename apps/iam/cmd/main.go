package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/golib/biz/gcontext"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/biz/gmiddleware/ginmiddleware"
	"github.com/morehao/golib/glog"
)

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

func main() {
	if err := serverInit(); err != nil {
		panic(fmt.Sprintf("server init failed, error: %v", err))
	}
	if config.Conf.Server.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	defer glog.Close()

	engine := gin.New()
	engine.ContextWithFallback = true
	engine.Use(gin.Recovery())
	engine.Use(ginmiddleware.AccessLog())
	engine.Use(ginmiddleware.JWTAuth(config.Conf.JWT.SignKey))
	engine.Use(ginContextToStdContext())
	iam.Routers(engine)

	if err := engine.Run(fmt.Sprintf(":%s", config.Conf.Server.Port)); err != nil {
		glog.Errorf(context.Background(), "%s run fail, port:%s", iam.AppName, config.Conf.Server.Port)
		panic(err)
	} else {
		glog.Infof(context.Background(), "%s run success, port:%s", iam.AppName, config.Conf.Server.Port)
	}
}
