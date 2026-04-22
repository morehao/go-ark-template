package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrapplication"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

// applicationRouter 初始化应用管理路由信息
func applicationRouter(groups *ginserver.RouterGroups) {
	applicationCtr := ctrapplication.NewApplicationCtr()

	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/application/create", applicationCtr.Create)
	v1RouterGroup.POST("/application/delete", applicationCtr.Delete)
	v1RouterGroup.POST("/application/update", applicationCtr.Update)
	v1RouterGroup.GET("/application/detail", applicationCtr.Detail)
	v1RouterGroup.POST("/application/pageList", applicationCtr.PageList)
}
