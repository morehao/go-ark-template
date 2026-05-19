package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrfaq"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func faqRouter(groups *ginserver.RouterGroups) {
	ctr := ctrfaq.NewFAQCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/faq/create", ctr.Create)
	v1RouterGroup.POST("/faq/delete", ctr.Delete)
	v1RouterGroup.POST("/faq/update", ctr.Update)
	v1RouterGroup.POST("/faq/pageList", ctr.PageList)
	v1RouterGroup.POST("/faq/search", ctr.Search)
	v1RouterGroup.GET("/faq/detail", ctr.Detail)
	v1RouterGroup.POST("/faq/import", ctr.Import)
}
