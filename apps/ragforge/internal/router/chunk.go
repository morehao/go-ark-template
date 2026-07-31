package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrchunk"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func chunkRouter(groups *ginserver.RouterGroups) {
	ctr := ctrchunk.NewChunkCtr()
	v1RouterGroup := groups.MustGetGroup(ginserver.ApiVersionV1)

	v1RouterGroup.POST("/chunk/delete", ctr.Delete)
	v1RouterGroup.POST("/chunk/update", ctr.Update)
	v1RouterGroup.POST("/chunk/pageList", ctr.PageList)
	v1RouterGroup.POST("/chunk/search", ctr.Search)
	v1RouterGroup.GET("/chunk/detail", ctr.Detail)
}
