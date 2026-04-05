package demo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demo/config"
	_ "github.com/morehao/goark/apps/demo/docs"
	"github.com/morehao/goark/apps/demo/internal/router"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

const AppName = "demo"

func Routers(engine *gin.Engine) {
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
