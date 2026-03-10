package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/demoapp"
	"github.com/morehao/goark/apps/demoapp/config"
	"github.com/morehao/golib/biz/gmiddleware/ginmiddleware"
	"github.com/morehao/golib/glog"
)

func main() {
	if err := serverInit(); err != nil {
		panic(fmt.Sprintf("server init failed, error: %v", err))
	}
	if config.Conf.Server.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	defer glog.Close()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(ginmiddleware.AccessLog())
	demoapp.Routers(engine)

	if err := engine.Run(fmt.Sprintf(":%s", config.Conf.Server.Port)); err != nil {
		glog.Errorf(context.Background(), "%s run fail, port:%s", demoapp.AppName, config.Conf.Server.Port)
		panic(err)
	} else {
		glog.Infof(context.Background(), "%s run success, port:%s", demoapp.AppName, config.Conf.Server.Port)
	}
}
