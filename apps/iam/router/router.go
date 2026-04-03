package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/config"
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
	iamGroup := v1AuthGroup.Group("/iam")
	// v1NoAuth := groups.NoAuthGroup.Group("/v1")
	organizationRouter(iamGroup)
	tenantRouter(iamGroup)
	departmentRouter(iamGroup)
	userRouter(iamGroup)
	menuRouter(iamGroup)
	roleRouter(iamGroup)
}
