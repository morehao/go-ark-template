package iam

import (
	"github.com/gin-gonic/gin"
	_ "github.com/morehao/goark/apps/iam/docs"
	"github.com/morehao/goark/apps/iam/internal/middleware"
	"github.com/morehao/goark/apps/iam/router"
)

const AppName = "iam"

func Routers(engine *gin.Engine, signKey string) {
	routerGroup := engine.Group(AppName)

	noAuthGroup := routerGroup.Group("")

	authGroup := routerGroup.Group("")
	authGroup.Use(middleware.JWTAuth(signKey))

	routerGroups := &router.RouterGroups{
		AuthGroup:   authGroup,
		NoAuthGroup: noAuthGroup,
	}
	router.RegisterRouter(routerGroups, AppName)
}
