package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrauth"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func authRouter(groups *ginrouter.RouterGroups) {
	authCtr := ctrauth.NewAuthCtr()

	groups.V1.POST("/auth/login", authCtr.Login)
	groups.V1.POST("/auth/selectTenant", authCtr.SelectTenant)
	groups.V1.POST("/auth/logout", authCtr.Logout)
}
