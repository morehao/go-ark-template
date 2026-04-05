package router

import (
	"github.com/morehao/goark/apps/demo/config"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func RegisterRouter(groups *ginrouter.RouterGroups, appName string) {
	if config.Conf.Server.Env == "dev" {
		ginrouter.RegisterSwagger(groups.V1, appName)
	}
	formatRouter(groups)
	sseRouter(groups)
	clientRouter(groups)
	userRouter(groups)
}
