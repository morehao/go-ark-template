package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrorg"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func organizationRouter(groups *ginserver.RouterGroups) {
	organizationCtr := ctrorg.NewOrganizationCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/organization/create", organizationCtr.Create)
	v1RouterGroup.POST("/organization/delete", organizationCtr.Delete)
	v1RouterGroup.POST("/organization/update", organizationCtr.Update)
	v1RouterGroup.GET("/organization/detail", organizationCtr.Detail)
	v1RouterGroup.POST("/organization/pageList", organizationCtr.PageList)
	v1RouterGroup.GET("/organization/getConfigsByDomain", organizationCtr.GetConfigsByDomain)
	v1RouterGroup.GET("/organization/listConfig", organizationCtr.ListConfig)
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
