package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctrauth"
)

func authRouter(authGroup, noAuthGroup *gin.RouterGroup) {
	authCtr := ctrauth.NewAuthCtr()

	// No-auth routes
	noAuthGroup.POST("/auth/login", authCtr.Login)
	noAuthGroup.POST("/auth/selectTenant", authCtr.SelectTenant)

	// Auth routes
	authGroup.POST("/auth/logout", authCtr.Logout)
	authGroup.POST("/auth/switchTenant", authCtr.SwitchTenant)
	authGroup.GET("/auth/currentUser", authCtr.GetCurrentUser)
}
