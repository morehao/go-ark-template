package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demo/internal/controller/ctruser"
)

// userRouter 初始化用户管理路由信息
func userRouter(routerGroup *gin.RouterGroup) {
	userCtr := ctruser.NewUserCtr()

	routerGroup.POST("/user/create", userCtr.Create)
	routerGroup.POST("/user/delete", userCtr.Delete)
	routerGroup.POST("/user/update", userCtr.Update)
	routerGroup.GET("/user/detail", userCtr.Detail)
	routerGroup.POST("/user/pageList", userCtr.PageList)
}
