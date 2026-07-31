package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrtag"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func tagRouter(groups *ginserver.RouterGroups) {
	ctr := ctrtag.NewTagCtr()
	v1RouterGroup := groups.MustGetGroup(ginserver.ApiVersionV1)

	v1RouterGroup.POST("/tag/create", ctr.Create)
	v1RouterGroup.POST("/tag/delete", ctr.Delete)
	v1RouterGroup.POST("/tag/update", ctr.Update)
	v1RouterGroup.GET("/tag/list", ctr.List)
}
