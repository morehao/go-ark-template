package router

import "github.com/morehao/golib/biz/gserver/ginserver"

func RegisterRouter(groups *ginserver.RouterGroups, appName string) {

	authRouter(groups)
	orgRouter(groups)
	tenantRouter(groups)
	departmentRouter(groups)
	userRouter(groups)
	menuRouter(groups)
	roleRouter(groups)
	applicationRouter(groups)
	apiKeyRouter(groups)
	oidcRouter(groups)
}
