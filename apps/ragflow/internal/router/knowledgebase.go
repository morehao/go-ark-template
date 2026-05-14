package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/ragflow/internal/controller/ctrknowledge"
)

func knowledgeBaseRouter(routerGroup *gin.RouterGroup) {
	knowledgeBaseCtr := ctrknowledge.NewKnowledgeBaseCtr()
	routerGroup.POST("/knowledgebase/create", knowledgeBaseCtr.Create)
	routerGroup.POST("/knowledgebase/delete", knowledgeBaseCtr.Delete)
	routerGroup.POST("/knowledgebase/update", knowledgeBaseCtr.Update)
	routerGroup.GET("/knowledgebase/detail", knowledgeBaseCtr.Detail)
	routerGroup.GET("/knowledgebase/list", knowledgeBaseCtr.List)
}