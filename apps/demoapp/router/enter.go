package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demoapp/config"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

type RouterGroups struct {
	AuthGroup   *gin.RouterGroup
	NoAuthGroup *gin.RouterGroup
}

func RegisterRouter(groups *RouterGroups, appName string) {
	if config.Conf.Server.Env == "dev" {
		ginrouter.RegisterSwagger(groups.AuthGroup, appName)
	}
	v1AuthGroup := groups.AuthGroup.Group("/v1")
	// v1NoAuth := groups.NoAuthGroup.Group("/v1")
	formatRouter(v1AuthGroup)
	sseRouter(v1AuthGroup)
	clientRouter(v1AuthGroup)
	userRouter(v1AuthGroup)
}
