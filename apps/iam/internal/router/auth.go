package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrauth"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func authRouter(groups *ginserver.RouterGroups) {
	authCtr := ctrauth.NewAuthCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/auth/loginByPassword", authCtr.LoginByPassword)
	v1RouterGroup.POST("/auth/selectTenant", authCtr.SelectTenant)
	v1RouterGroup.POST("/auth/refreshToken", authCtr.RefreshToken)
	v1RouterGroup.POST("/auth/logout", authCtr.Logout)
}
