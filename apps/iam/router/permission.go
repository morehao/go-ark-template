package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctrpermission"
)

func menuRouter(routerGroup *gin.RouterGroup) {
	menuCtr := ctrpermission.NewMenuCtr()

	routerGroup.POST("/menu/create", menuCtr.Create)
	routerGroup.POST("/menu/delete", menuCtr.Delete)
	routerGroup.POST("/menu/update", menuCtr.Update)
	routerGroup.GET("/menu/detail", menuCtr.Detail)
	routerGroup.POST("/menu/pageList", menuCtr.PageList)
}

func roleRouter(routerGroup *gin.RouterGroup) {
	roleCtr := ctrpermission.NewRoleCtr()
	routerGroup.POST("/role/create", roleCtr.Create)
	routerGroup.POST("/role/delete", roleCtr.Delete)
	routerGroup.POST("/role/update", roleCtr.Update)
	routerGroup.GET("/role/detail", roleCtr.Detail)
	routerGroup.POST("/role/pageList", roleCtr.PageList)
	routerGroup.POST("/role/assignMenus", roleCtr.AssignMenus)
	routerGroup.POST("/role/removeMenus", roleCtr.RemoveMenus)
	routerGroup.GET("/role/listMenus", roleCtr.ListMenus)

	authorizationCtr := ctrpermission.NewAuthorizationCtr()
	routerGroup.POST("/user/assignRoles", authorizationCtr.AssignRoles)
	routerGroup.POST("/user/removeRoles", authorizationCtr.RemoveRoles)
	routerGroup.GET("/user/listRoles", authorizationCtr.ListUserRoles)
	routerGroup.GET("/permission/userPermissions", authorizationCtr.GetUserPermissions)
}
