package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctroidc"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func oidcRouter(groups *ginserver.RouterGroups) {
	authorizeCtr := ctroidc.NewAuthorizeCtr()
	tokenCtr := ctroidc.NewTokenCtr()
	userinfoCtr := ctroidc.NewUserinfoCtr()
	logoutCtr := ctroidc.NewLogoutCtr()

	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.GET("/oidc/authorize", authorizeCtr.Authorize)
	v1RouterGroup.POST("/oidc/token", tokenCtr.Token)
	v1RouterGroup.POST("/oidc/refreshToken", tokenCtr.RefreshToken)
	v1RouterGroup.GET("/oidc/userinfo", userinfoCtr.UserInfo)
	v1RouterGroup.POST("/oidc/logout", logoutCtr.Logout)
}