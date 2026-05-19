package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrmodel"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func modelRouter(groups *ginserver.RouterGroups) {
	ctr := ctrmodel.NewModelCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/model/create", ctr.Create)
	v1RouterGroup.POST("/model/delete", ctr.Delete)
	v1RouterGroup.POST("/model/update", ctr.Update)
	v1RouterGroup.GET("/model/detail", ctr.Detail)
	v1RouterGroup.POST("/model/pageList", ctr.PageList)
	v1RouterGroup.POST("/model/test", ctr.Test)
	v1RouterGroup.GET("/model/providers", ctr.GetProviders)
}
