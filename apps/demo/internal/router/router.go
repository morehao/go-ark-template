package router

import "github.com/morehao/golib/biz/gserver/ginserver"

func RegisterRouter(groups *ginserver.RouterGroups, appName string) {
	formatRouter(groups)
	sseRouter(groups)
	clientRouter(groups)
	userRouter(groups)
}
