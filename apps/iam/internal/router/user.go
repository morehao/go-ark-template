package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctruser"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func userRouter(groups *ginrouter.RouterGroups) {
	userCtr := ctruser.NewUserCtr()

	groups.V1.POST("/user/create", userCtr.Create)
	groups.V1.POST("/user/delete", userCtr.Delete)
	groups.V1.POST("/user/update", userCtr.Update)
	groups.V1.GET("/user/detail", userCtr.Detail)
	groups.V1.POST("/user/pageList", userCtr.PageList)
	groups.V1.POST("/user/assignDepartment", userCtr.AssignDepartment)
	groups.V1.POST("/user/removeDepartment", userCtr.RemoveDepartment)
	groups.V1.GET("/user/listDepartments", userCtr.ListDepartments)
	groups.V1.POST("/user/assignRoles", userCtr.AssignRoles)
	groups.V1.POST("/user/removeRoles", userCtr.RemoveRoles)
	groups.V1.GET("/user/listRoles", userCtr.ListRoles)
}
