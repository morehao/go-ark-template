package iam

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
	_ "github.com/morehao/goark/apps/iam/docs"
	"github.com/morehao/goark/apps/iam/internal/router"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gmiddleware/ginmiddleware"
	"github.com/morehao/golib/biz/gserver/gindocs"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

const AppName = "iam"

func Routers(engine *gin.Engine) {
	routerGroups := ginserver.NewRouterGroups(engine, AppName, ginserver.Version{
		Name: gconstant.ApiVersionV1,
		Middlewares: []gin.HandlerFunc{
			ginmiddleware.JWTAuth(config.Conf.JWT.SignKey, ginmiddleware.WithAuthSkipPaths(
				"/v1/iam/org/getConfigsByDomain",
				"/v1/iam/auth/register",
				"/v1/iam/auth/loginByPassword",
				"/v1/iam/auth/selectTenant",
			)),
			ginmiddleware.TokenBlacklistCheck(dbclient.RedisCli, ginmiddleware.WithBlacklistKeyPrefix("iam:token:blacklist:")),
		},
	})

	if config.Conf.Server.Env == "dev" {
		gindocs.Register(engine.Group("/"+AppName), AppName)
	}

	router.RegisterRouter(routerGroups, AppName)
}
