package router

import (
	"github.com/morehao/goark/apps/iam/internal/controller/ctrlog"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func logRouter(groups *ginserver.RouterGroups) {
	operationLogCtr := ctrlog.NewOperationLogCtr()
	loginLogCtr := ctrlog.NewLoginLogCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)

	v1RouterGroup.POST("/iam/operationLog/create", operationLogCtr.Create)
	v1RouterGroup.GET("/iam/operationLog/pageList", operationLogCtr.PageList)

	v1RouterGroup.POST("/iam/loginLog/create", loginLogCtr.Create)
	v1RouterGroup.GET("/iam/loginLog/pageList", loginLogCtr.PageList)
}
