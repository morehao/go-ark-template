package main

import (
	"context"
	"fmt"

	"github.com/morehao/goark/apps/ragflow"
	"github.com/morehao/goark/apps/ragflow/config"
	"github.com/morehao/golib/glog"
)

func serverInit() error {
	config.InitConf()
	defaultLogCfg := config.Conf.Log["default"]
	if err := glog.InitLogger(&defaultLogCfg); err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}
	return nil
}

func shutdownTraceProvider() {
	glog.Close()
}