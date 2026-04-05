package router

import (
	"github.com/morehao/goark/apps/demo/internal/controller/ctrexample"
	"github.com/morehao/golib/biz/grouter/ginrouter"
)

func formatRouter(groups *ginrouter.RouterGroups) {
	formatCtr := ctrexample.NewFormatCtr()
	groups.V1.GET("/formatRes", formatCtr.FormatRes)
}

func sseRouter(groups *ginrouter.RouterGroups) {
	sseCtr := ctrexample.NewSSECtr()
	groups.V1.GET("/time", sseCtr.Time)
	groups.V1.GET("/timeRaw", sseCtr.TimeRaw)
	groups.V1.GET("/process", sseCtr.Process)
	groups.V1.GET("/chat", sseCtr.Chat)
	groups.V1.GET("/raw", sseCtr.Raw)
}

func clientRouter(groups *ginrouter.RouterGroups) {
	clientCtr := ctrexample.NewClientCtr()
	groups.V1.GET("/CallGetHttpbingo", clientCtr.CallGetHttpbingo)
}
