package demoapp

import (
	"github.com/gin-gonic/gin"
	_ "github.com/morehao/goark/apps/demoapp/docs"
	"github.com/morehao/goark/apps/demoapp/router"
)

const AppName = "demoapp"

func Routers(engine *gin.Engine) {
	routerGroup := engine.Group(AppName)
	routerGroups := &router.RouterGroups{
		AuthGroup:   routerGroup,
		NoAuthGroup: routerGroup,
	}
	router.RegisterRouter(routerGroups, AppName)
}
