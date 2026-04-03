package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctruser"
)

func userRouter(routerGroup *gin.RouterGroup) {
	userCtr := ctruser.NewUserCtr()

	routerGroup.POST("/user/create", userCtr.Create)
	routerGroup.POST("/user/delete", userCtr.Delete)
	routerGroup.POST("/user/update", userCtr.Update)
	routerGroup.GET("/user/detail", userCtr.Detail)
	routerGroup.POST("/user/pageList", userCtr.PageList)
	routerGroup.POST("/user/assignDepartment", userCtr.AssignDepartment)
	routerGroup.POST("/user/removeDepartment", userCtr.RemoveDepartment)
	routerGroup.GET("/user/listDepartments", userCtr.ListDepartments)
	routerGroup.POST("/user/assignRoles", userCtr.AssignRoles)
	routerGroup.POST("/user/removeRoles", userCtr.RemoveRoles)
	routerGroup.GET("/user/listRoles", userCtr.ListRoles)
}
