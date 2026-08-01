package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrkb"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func kbRouter(groups *ginserver.RouterGroups) {
	ctr := ctrkb.NewKBCtr()
	v1RouterGroup := groups.MustGetGroup(ginserver.ApiVersionV1)

	v1RouterGroup.POST("/kb/create", ctr.Create)
	v1RouterGroup.POST("/kb/delete", ctr.Delete)
	v1RouterGroup.POST("/kb/update", ctr.Update)
	v1RouterGroup.GET("/kb/detail", ctr.Detail)
	v1RouterGroup.POST("/kb/pageList", ctr.PageList)
	v1RouterGroup.POST("/kb/copy", ctr.Copy)
}
