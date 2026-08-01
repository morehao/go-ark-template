package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrmessage"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func messageRouter(groups *ginserver.RouterGroups) {
	ctr := ctrmessage.NewMessageCtr()
	v1RouterGroup := groups.MustGetGroup(ginserver.ApiVersionV1)

	v1RouterGroup.GET("/message/list", ctr.List)
	v1RouterGroup.POST("/message/delete", ctr.Delete)
	v1RouterGroup.POST("/message/search", ctr.Search)
}
