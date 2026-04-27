package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctroidc"
)

func oidcRouter(routerGroup *gin.RouterGroup) {
	oidcGroup := routerGroup.Group("/oidc")

	authorizeCtr := ctroidc.NewAuthorizeCtr()
	tokenCtr := ctroidc.NewTokenCtr()
	userinfoCtr := ctroidc.NewUserinfoCtr()
	logoutCtr := ctroidc.NewLogoutCtr()

	oidcGroup.GET("/authorize", authorizeCtr.Authorize)
	oidcGroup.POST("/token", tokenCtr.Token)
	oidcGroup.POST("/token/refresh", tokenCtr.RefreshToken)
	oidcGroup.GET("/userinfo", userinfoCtr.UserInfo)
	oidcGroup.POST("/logout", logoutCtr.Logout)
}