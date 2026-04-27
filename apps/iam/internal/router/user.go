package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctruser"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func userRouter(groups *ginserver.RouterGroups) {
	userCtr := ctruser.NewUserCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/user/create", userCtr.Create)
	v1RouterGroup.POST("/user/delete", userCtr.Delete)
	v1RouterGroup.POST("/user/update", userCtr.Update)
	v1RouterGroup.GET("/user/detail", userCtr.Detail)
	v1RouterGroup.POST("/user/pageList", userCtr.PageList)
	v1RouterGroup.POST("/user/assignDepartments", userCtr.AssignDepartments)
	v1RouterGroup.GET("/user/listDepartments", userCtr.ListDepartments)
	v1RouterGroup.POST("/user/assignRoles", userCtr.AssignRoles)
	v1RouterGroup.GET("/user/listRoles", userCtr.ListRoles)
	v1RouterGroup.GET("/user/getCurrentUserInfo", userCtr.GetCurrentUserInfo)
}
