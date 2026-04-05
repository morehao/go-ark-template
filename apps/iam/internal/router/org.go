package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrorg"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func tenantRouter(groups *ginrouter.RouterGroups) {
	tenantCtr := ctrorg.NewTenantCtr()

	groups.V1.POST("/tenant/create", tenantCtr.Create)
	groups.V1.POST("/tenant/delete", tenantCtr.Delete)
	groups.V1.POST("/tenant/update", tenantCtr.Update)
	groups.V1.GET("/tenant/detail", tenantCtr.Detail)
	groups.V1.POST("/tenant/pageList", tenantCtr.PageList)
}

func departmentRouter(groups *ginrouter.RouterGroups) {
	departmentCtr := ctrorg.NewDepartmentCtr()
	groups.V1.POST("/department/create", departmentCtr.Create)
	groups.V1.POST("/department/delete", departmentCtr.Delete)
	groups.V1.POST("/department/update", departmentCtr.Update)
	groups.V1.GET("/department/detail", departmentCtr.Detail)
	groups.V1.POST("/department/pageList", departmentCtr.PageList)
	groups.V1.GET("/department/tree", departmentCtr.Tree)
}
