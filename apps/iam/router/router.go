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
	iamAuthGroup := v1AuthGroup.Group("/iam")

	v1NoAuthGroup := groups.NoAuthGroup.Group("/v1")
	iamNoAuthGroup := v1NoAuthGroup.Group("/iam")

	// 认证路由(登录在NoAuth组，选择租户和登出在Auth组)
	authRouter(iamNoAuthGroup, iamAuthGroup)

	// 需要认证的路由
	organizationRouter(iamAuthGroup)
	tenantRouter(iamAuthGroup)
	departmentRouter(iamAuthGroup)
	userRouter(iamAuthGroup)
	menuRouter(iamAuthGroup)
	roleRouter(iamAuthGroup)
}
