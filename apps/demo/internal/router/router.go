package router

import "github.com/morehao/golib/biz/gserver/ginserver"

func RegisterRouter(groups *ginserver.RouterGroups) {
	formatRouter(groups)
	sseRouter(groups)
	clientRouter(groups)
	userRouter(groups)
}
