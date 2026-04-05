package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrorg"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func organizationRouter(groups *ginrouter.RouterGroups) {
	organizationCtr := ctrorg.NewOrganizationCtr()

	groups.V1.POST("/organization/create", organizationCtr.Create)
	groups.V1.POST("/organization/delete", organizationCtr.Delete)
	groups.V1.POST("/organization/update", organizationCtr.Update)
	groups.V1.GET("/organization/detail", organizationCtr.Detail)
	groups.V1.POST("/organization/pageList", organizationCtr.PageList)
}
