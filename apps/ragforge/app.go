package ragforge

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/config"
	_ "github.com/morehao/goark/ragforge/docs"
	"github.com/morehao/goark/ragforge/internal/middleware"
	"github.com/morehao/goark/ragforge/internal/router"
	"github.com/morehao/golib/biz/gserver/gindocs"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

const AppName = "ragforge"

func Routers(engine *gin.Engine) {
	routerGroups := ginserver.NewRouterGroups(engine, AppName, ginserver.VersionGroup{
		Version: ginserver.ApiVersionV1,
		Middlewares: []gin.HandlerFunc{
			middleware.Auth(),
		},
	})

	if config.Conf.Server.Env == "dev" {
		gindocs.Register(engine.Group("/"+AppName), AppName)
	}
	router.RegisterRouter(routerGroups, AppName)
}
