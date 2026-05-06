package router

import (
	"github.com/morehao/goark/iam/internal/controller/ctrorg"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func orgRouter(groups *ginserver.RouterGroups) {
	orgCtr := ctrorg.NewOrganizationCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/organization/create", orgCtr.Create)
	v1RouterGroup.POST("/organization/delete", orgCtr.Delete)
	v1RouterGroup.POST("/organization/update", orgCtr.Update)
	v1RouterGroup.GET("/organization/detail", orgCtr.Detail)
	v1RouterGroup.POST("/organization/pageList", orgCtr.PageList)
	v1RouterGroup.GET("/organization/getOrgConfig", orgCtr.GetOrgConfig)
	v1RouterGroup.GET("/organization/listConfigDefinitions", orgCtr.ListConfigDefinitions)
}

func tenantRouter(groups *ginserver.RouterGroups) {
	tenantCtr := ctrorg.NewTenantCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/tenant/create", tenantCtr.Create)
	v1RouterGroup.POST("/tenant/delete", tenantCtr.Delete)
	v1RouterGroup.POST("/tenant/update", tenantCtr.Update)
	v1RouterGroup.GET("/tenant/detail", tenantCtr.Detail)
	v1RouterGroup.POST("/tenant/pageList", tenantCtr.PageList)
}

func departmentRouter(groups *ginserver.RouterGroups) {
	departmentCtr := ctrorg.NewDepartmentCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)
	v1RouterGroup.POST("/department/create", departmentCtr.Create)
	v1RouterGroup.POST("/department/delete", departmentCtr.Delete)
	v1RouterGroup.POST("/department/update", departmentCtr.Update)
	v1RouterGroup.GET("/department/detail", departmentCtr.Detail)
	v1RouterGroup.POST("/department/pageList", departmentCtr.PageList)
	v1RouterGroup.GET("/department/tree", departmentCtr.Tree)
}
