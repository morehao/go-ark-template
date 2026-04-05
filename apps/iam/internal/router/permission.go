package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrpermission"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func menuRouter(groups *ginrouter.RouterGroups) {
	menuCtr := ctrpermission.NewMenuCtr()

	groups.V1.POST("/menu/create", menuCtr.Create)
	groups.V1.POST("/menu/delete", menuCtr.Delete)
	groups.V1.POST("/menu/update", menuCtr.Update)
	groups.V1.GET("/menu/detail", menuCtr.Detail)
	groups.V1.POST("/menu/pageList", menuCtr.PageList)
	groups.V1.GET("/menu/tree", menuCtr.Tree)
}

func roleRouter(groups *ginrouter.RouterGroups) {
	roleCtr := ctrpermission.NewRoleCtr()
	groups.V1.POST("/role/create", roleCtr.Create)
	groups.V1.POST("/role/delete", roleCtr.Delete)
	groups.V1.POST("/role/update", roleCtr.Update)
	groups.V1.GET("/role/detail", roleCtr.Detail)
	groups.V1.POST("/role/pageList", roleCtr.PageList)
	groups.V1.POST("/role/assignMenus", roleCtr.AssignMenus)
	groups.V1.GET("/role/listMenus", roleCtr.ListMenus)
}
