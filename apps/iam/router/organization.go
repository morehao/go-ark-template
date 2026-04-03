package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctrorg"
)

func organizationRouter(routerGroup *gin.RouterGroup) {
	organizationCtr := ctrorg.NewOrganizationCtr()

	routerGroup.POST("/organization/create", organizationCtr.Create)
	routerGroup.POST("/organization/delete", organizationCtr.Delete)
	routerGroup.POST("/organization/update", organizationCtr.Update)
	routerGroup.GET("/organization/detail", organizationCtr.Detail)
	routerGroup.POST("/organization/pageList", organizationCtr.PageList)
}
