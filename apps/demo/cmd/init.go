package main

import (
	"fmt"

	"github.com/morehao/goark/apps/demo/config"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/glog"
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
		return fmt.Errorf("init logger failed: " + err.Error())
	}
	return nil
}

func resourceInit() error {
	var gormLogConfig *glog.LogConfig
	if cfg, ok := config.Conf.Log["gorm"]; ok {
		gormLogConfig = &cfg
	}
	if err := dbclient.InitMultiDB(config.Conf.DBConfigs, gormLogConfig); err != nil {
		return fmt.Errorf("init db failed: " + err.Error())
	}
	var redisLogConfig *glog.LogConfig
	if cfg, ok := config.Conf.Log["redis"]; ok {
		redisLogConfig = &cfg
	}
	if err := dbclient.InitRedis(config.Conf.RedisConfig, redisLogConfig); err != nil {
		return fmt.Errorf("init redis failed: " + err.Error())
	}
	var esLogConfig *glog.LogConfig
	if cfg, ok := config.Conf.Log["es"]; ok {
		esLogConfig = &cfg
	}
	if err := dbclient.InitMultiEs(config.Conf.ESConfigs, esLogConfig); err != nil {
		return fmt.Errorf("init es failed: " + err.Error())
	}
	return nil
}
