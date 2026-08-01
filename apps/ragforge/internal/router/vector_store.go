package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrvectorstore"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func vectorStoreRouter(groups *ginserver.RouterGroups) {
	ctr := ctrvectorstore.NewVectorStoreCtr()
	v1RouterGroup := groups.MustGetGroup(ginserver.ApiVersionV1)

	v1RouterGroup.POST("/vectorStore/create", ctr.Create)
	v1RouterGroup.POST("/vectorStore/delete", ctr.Delete)
	v1RouterGroup.POST("/vectorStore/update", ctr.Update)
	v1RouterGroup.GET("/vectorStore/detail", ctr.Detail)
	v1RouterGroup.POST("/vectorStore/pageList", ctr.PageList)
	v1RouterGroup.POST("/vectorStore/test", ctr.Test)
	v1RouterGroup.GET("/vectorStore/types", ctr.GetTypes)
}
