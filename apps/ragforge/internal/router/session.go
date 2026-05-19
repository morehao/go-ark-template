package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrsession"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func sessionRouter(groups *ginserver.RouterGroups) {
	ctr := ctrsession.NewSessionCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/session/create", ctr.Create)
	v1RouterGroup.POST("/session/delete", ctr.Delete)
	v1RouterGroup.POST("/session/update", ctr.Update)
	v1RouterGroup.GET("/session/detail", ctr.Detail)
	v1RouterGroup.POST("/session/pageList", ctr.PageList)
	v1RouterGroup.POST("/session/generateTitle", ctr.GenerateTitle)
	v1RouterGroup.POST("/session/stop", ctr.Stop)
}
