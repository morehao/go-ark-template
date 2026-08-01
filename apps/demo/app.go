package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/demo/config"
	_ "github.com/morehao/goark/demo/docs"
	"github.com/morehao/goark/demo/internal/middleware"
	"github.com/morehao/goark/demo/internal/router"
	"github.com/morehao/golib/biz/gserver/gindocs"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

const AppName = "demo"

func Routers(engine *gin.Engine) {
	routerGroups := ginserver.NewRouterGroups(engine, AppName, ginserver.VersionGroup{
		Version: ginserver.ApiVersionV1,
		Middlewares: []gin.HandlerFunc{
			middleware.Example(),
		},
	})

	if config.Conf.Server.Env == "dev" {
		gindocs.Register(engine.Group("/"+AppName), AppName)
	}

	router.RegisterRouter(routerGroups, AppName)
}
