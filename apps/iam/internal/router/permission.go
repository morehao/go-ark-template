package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrpermission"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func menuRouter(groups *ginserver.RouterGroups) {
	menuCtr := ctrpermission.NewMenuCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/menu/create", menuCtr.Create)
	v1RouterGroup.POST("/menu/delete", menuCtr.Delete)
	v1RouterGroup.POST("/menu/update", menuCtr.Update)
	v1RouterGroup.GET("/menu/detail", menuCtr.Detail)
	v1RouterGroup.POST("/menu/pageList", menuCtr.PageList)
	v1RouterGroup.GET("/menu/tree", menuCtr.Tree)
}

func roleRouter(groups *ginserver.RouterGroups) {
	roleCtr := ctrpermission.NewRoleCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)
	v1RouterGroup.POST("/role/create", roleCtr.Create)
	v1RouterGroup.POST("/role/delete", roleCtr.Delete)
	v1RouterGroup.POST("/role/update", roleCtr.Update)
	v1RouterGroup.GET("/role/detail", roleCtr.Detail)
	v1RouterGroup.POST("/role/pageList", roleCtr.PageList)
	v1RouterGroup.POST("/role/assignMenus", roleCtr.AssignMenus)
	v1RouterGroup.GET("/role/listMenus", roleCtr.ListMenus)
}
