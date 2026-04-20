package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/morehao/goark/apps/iam"
	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gtrace"
	"github.com/morehao/golib/gtrace/otlptracegrpc"
)

var traceProvider *gtrace.Provider

func serverInit() error {
	if err := preInit(); err != nil {
		return err
	}
	if err := initTrace(); err != nil {
		return err
	}
	if err := resourceInit(); err != nil {
		return err
	}
	return nil
}

func preInit() error {
	config.InitConf()
	defaultLogCfg := config.Conf.Log["default"]
	if err := glog.InitLogger(&defaultLogCfg); err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}
	return nil
}

func resourceInit() error {
	var gormLogConfig *glog.LogConfig
	if cfg, ok := config.Conf.Log["gorm"]; ok {
		gormLogConfig = &cfg
	}
	if err := dbclient.InitMultiDB(config.Conf.DBConfigs, gormLogConfig); err != nil {
		return fmt.Errorf("init db failed: %w", err)
	}
	
	var redisLogConfig *glog.LogConfig
	if cfg, ok := config.Conf.Log["redis"]; ok {
		redisLogConfig = &cfg
	}
	if err := dbclient.InitRedis(config.Conf.RedisConfig, redisLogConfig); err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}
	return nil
}

func initTrace() error {
	traceCfg := config.Conf.Trace
	if !traceCfg.Enable {
		glog.Infof(context.Background(), "[%s.initTrace] trace disabled, skip init", iam.AppName)
		return nil
	}
	if strings.TrimSpace(traceCfg.OTLP.Endpoint) == "" {
		glog.Infof(context.Background(), "[%s.initTrace] trace enabled but otlp endpoint empty, skip init", iam.AppName)
		return nil
	}

	tCfg := gtrace.DefaultConfig(iam.AppName)
	tCfg.ServiceVersion = traceCfg.ServiceVersion
	tCfg.Environment = config.Conf.Server.Env
	tCfg.TraceIDRatio = traceCfg.TraceIDRatio
	sampler, err := gtrace.ParseSampler(traceCfg.Sampler)
	if err != nil {
		return fmt.Errorf("init trace failed: %w", err)
	}
	tCfg.Sampler = sampler

	eCfg := otlptracegrpc.DefaultConfig()
	eCfg.Endpoint = traceCfg.OTLP.Endpoint
	eCfg.Insecure = traceCfg.OTLP.Insecure
	if traceCfg.OTLP.Timeout > 0 {
		eCfg.Timeout = traceCfg.OTLP.Timeout
	}

	provider, err := gtrace.Init(context.Background(), tCfg, otlptracegrpc.NewExporterFactory(eCfg))
	if err != nil {
		glog.Errorf(context.Background(), "[%s.initTrace] init trace failed, fallback to disabled mode, err:%v", iam.AppName, err)
		return nil
	}
	traceProvider = provider
	return nil
}

func shutdownTraceProvider() {
	if traceProvider == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := traceProvider.Shutdown(ctx); err != nil {
		glog.Errorf(context.Background(), "[%s.shutdownTraceProvider] shutdown fail, err:%v", iam.AppName, err)
	}
}
