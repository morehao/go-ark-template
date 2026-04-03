package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demo/config"
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
	formatRouter(groups.AuthGroup)
	sseRouter(groups.AuthGroup)
	clientRouter(groups.AuthGroup)
	userRouter(groups.AuthGroup)
}
