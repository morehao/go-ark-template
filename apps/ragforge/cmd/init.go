package main

import (
	"context"
	"fmt"

	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/config"
	"github.com/morehao/goark/ragforge/internal/database"
	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/internal/retriever"
	"github.com/morehao/golib/glog"
	_ "github.com/morehao/golib/glog/driver/slog"
)

func serverInit() error {
	if err := preInit(); err != nil {
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
	if err := database.RunMigrations(context.Background()); err != nil {
		glog.Errorf(context.Background(), "[main.resourceInit] run migrations fail, err:%v", err)
	}
	var redisLogConfig *glog.LogConfig
	if cfg, ok := config.Conf.Log["redis"]; ok {
		redisLogConfig = &cfg
	}
	if err := dbclient.InitRedis(config.Conf.RedisConfig, redisLogConfig); err != nil {
		return fmt.Errorf("init redis failed: %w", err)
	}

	engine.NewEngineFactory()
	retriever.Init(retriever.NewRetriever())
	glog.Infof(context.Background(), "[main.resourceInit] engine factory and retriever initialized")

	return nil
}
