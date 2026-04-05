package router

import "github.com/morehao/golib/biz/grouter/ginrouter"

func RegisterRouter(groups *ginrouter.RouterGroups, appName string) {

	authRouter(groups)
	organizationRouter(groups)
	tenantRouter(groups)
	departmentRouter(groups)
	userRouter(groups)
	menuRouter(groups)
	roleRouter(groups)
}
