package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrknowledge"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func knowledgeRouter(groups *ginserver.RouterGroups) {
	ctr := ctrknowledge.NewKnowledgeCtr()
	v1RouterGroup := groups.MustGetGroup(ginserver.ApiVersionV1)

	v1RouterGroup.POST("/knowledge/createFile", ctr.CreateFromFile)
	v1RouterGroup.POST("/knowledge/createUrl", ctr.CreateFromURL)
	v1RouterGroup.POST("/knowledge/createManual", ctr.CreateManual)
	v1RouterGroup.POST("/knowledge/delete", ctr.Delete)
	v1RouterGroup.POST("/knowledge/update", ctr.Update)
	v1RouterGroup.POST("/knowledge/reparse", ctr.Reparse)
	v1RouterGroup.POST("/knowledge/pageList", ctr.PageList)
	v1RouterGroup.GET("/knowledge/detail", ctr.Detail)
	v1RouterGroup.GET("/knowledge/download", ctr.Download)
}
