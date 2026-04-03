package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctrauth"
)

func authRouter(noAuthGroup *gin.RouterGroup, authGroup *gin.RouterGroup) {
	authCtr := ctrauth.NewAuthCtr()

	// 登录接口不需要认证
	noAuthGroup.POST("/auth/login", authCtr.Login)

	// 选择租户和登出需要认证
	authGroup.POST("/auth/selectTenant", authCtr.SelectTenant)
	authGroup.POST("/auth/logout", authCtr.Logout)
}
