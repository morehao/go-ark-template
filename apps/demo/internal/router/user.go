package router

import (
	"github.com/morehao/goark/apps/demo/internal/controller/ctruser"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

// userRouter 初始化用户管理路由信息
func userRouter(groups *ginrouter.RouterGroups) {
	userCtr := ctruser.NewUserCtr()

	groups.V1.POST("/user/create", userCtr.Create)
	groups.V1.POST("/user/delete", userCtr.Delete)
	groups.V1.POST("/user/update", userCtr.Update)
	groups.V1.GET("/user/detail", userCtr.Detail)
	groups.V1.POST("/user/pageList", userCtr.PageList)
}
