package router

import (
	"github.com/morehao/goark/iam/internal/controller/ctrkey"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func apiKeyRouter(groups *ginserver.RouterGroups) {
	apiKeyCtr := ctrkey.NewApiKeyCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/iam/apiKey/create", apiKeyCtr.Create)
	v1RouterGroup.POST("/iam/apiKey/delete", apiKeyCtr.Delete)
	v1RouterGroup.GET("/iam/apiKey/list", apiKeyCtr.List)
	v1RouterGroup.POST("/iam/apiKey/disable", apiKeyCtr.Disable)
	v1RouterGroup.POST("/iam/apiKey/enable", apiKeyCtr.Enable)
}
