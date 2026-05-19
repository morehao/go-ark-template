package router

import (
	"github.com/morehao/goark/ragforge/internal/controller/ctrqa"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func qaRouter(groups *ginserver.RouterGroups) {
	ctr := ctrqa.NewQACtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/qa/knowledgeChat", ctr.KnowledgeChat)
	v1RouterGroup.POST("/qa/knowledgeChatStream", ctr.KnowledgeChatStream)
}
