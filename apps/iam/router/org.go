package router

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/controller/ctrorg"
)

func tenantRouter(routerGroup *gin.RouterGroup) {
	tenantCtr := ctrorg.NewTenantCtr()

	routerGroup.POST("/tenant/create", tenantCtr.Create)
	routerGroup.POST("/tenant/delete", tenantCtr.Delete)
	routerGroup.POST("/tenant/update", tenantCtr.Update)
	routerGroup.GET("/tenant/detail", tenantCtr.Detail)
	routerGroup.POST("/tenant/pageList", tenantCtr.PageList)
}

func departmentRouter(routerGroup *gin.RouterGroup) {
	departmentCtr := ctrorg.NewDepartmentCtr()
	routerGroup.POST("/department/create", departmentCtr.Create)
	routerGroup.POST("/department/delete", departmentCtr.Delete)
	routerGroup.POST("/department/update", departmentCtr.Update)
	routerGroup.GET("/department/detail", departmentCtr.Detail)
	routerGroup.POST("/department/pageList", departmentCtr.PageList)
	routerGroup.GET("/department/tree", departmentCtr.Tree)
}
