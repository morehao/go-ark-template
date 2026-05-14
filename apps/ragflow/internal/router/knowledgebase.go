package router

import (
	"github.com/morehao/goark/apps/ragflow/internal/controller/ctrknowledge"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func knowledgeBaseRouter(groups *ginserver.RouterGroups) {
	knowledgeBaseCtr := ctrknowledge.NewKnowledgeBaseCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/ragflow/knowledgebase/create", knowledgeBaseCtr.Create)
	v1RouterGroup.POST("/ragflow/knowledgebase/delete", knowledgeBaseCtr.Delete)
	v1RouterGroup.POST("/ragflow/knowledgebase/update", knowledgeBaseCtr.Update)
	v1RouterGroup.GET("/ragflow/knowledgebase/detail", knowledgeBaseCtr.Detail)
	v1RouterGroup.GET("/ragflow/knowledgebase/list", knowledgeBaseCtr.List)
}