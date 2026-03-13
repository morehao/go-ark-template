package demo

import (
	"github.com/gin-gonic/gin"
	_ "github.com/morehao/goark/apps/demo/docs"
	"github.com/morehao/goark/apps/demo/router"
)

const AppName = "demo"

func Routers(engine *gin.Engine) {
	routerGroup := engine.Group(AppName)
	routerGroups := &router.RouterGroups{
		AuthGroup:   routerGroup,
		NoAuthGroup: routerGroup,
	}
	router.RegisterRouter(routerGroups, AppName)
}
