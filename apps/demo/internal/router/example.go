package router

import (
	"github.com/morehao/goark/apps/demo/internal/controller/ctrexample"
	"github.com/morehao/golib/biz/gconstant"
	"github.com/morehao/golib/biz/gserver/ginserver"
)

func formatRouter(groups *ginserver.RouterGroups) {
	formatCtr := ctrexample.NewFormatCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)
	v1RouterGroup.GET("/formatRes", formatCtr.FormatRes)
}

func sseRouter(groups *ginserver.RouterGroups) {
	sseCtr := ctrexample.NewSSECtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)
	v1RouterGroup.GET("/time", sseCtr.Time)
	v1RouterGroup.GET("/timeRaw", sseCtr.TimeRaw)
	v1RouterGroup.GET("/process", sseCtr.Process)
	v1RouterGroup.GET("/chat", sseCtr.Chat)
	v1RouterGroup.GET("/raw", sseCtr.Raw)
}

func clientRouter(groups *ginserver.RouterGroups) {
	clientCtr := ctrexample.NewClientCtr()
	v1RouterGroup := groups.MustGetGroup(gconstant.ApiVersionV1)
	v1RouterGroup.GET("/CallGetHttpbingo", clientCtr.CallGetHttpbingo)
}
