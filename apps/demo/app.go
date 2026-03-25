package demo

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/morehao/goark/apps/demo/docs"
	"github.com/morehao/goark/apps/demo/router"
	"github.com/morehao/golib/biz/gconstant"
)

const AppName = "demo"

func Routers(engine *gin.Engine) {
	routerGroup := engine.Group(fmt.Sprintf("/%s/%s", gconstant.ApiVersionV1, AppName))
	routerGroups := &router.RouterGroups{
		AuthGroup:   routerGroup,
		NoAuthGroup: routerGroup,
	}
	router.RegisterRouter(routerGroups, AppName)
}
