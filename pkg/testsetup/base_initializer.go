package testsetup

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/dbaccess/dbes"
	"github.com/morehao/golib/dbaccess/dbgorm"
	"github.com/morehao/golib/dbaccess/dbredis"
	"github.com/morehao/golib/glog"
)

type baseAppInitializer struct {
	AppName    string
	ConfigPath string

	Log         map[string]glog.LogConfig
	DBConfigs   []dbgorm.GormConfig
	RedisConfig dbredis.RedisConfig
	ESConfigs   []dbes.ESConfig
}

func newBaseAppInitializer(appName string) (*baseAppInitializer, error) {
	configPath := findConfigPath(appName)
	if configPath == "" {
		return nil, fmt.Errorf("cannot find config path for app: %s", appName)
	}

	return &baseAppInitializer{
		AppName:    appName,
		ConfigPath: configPath,
	}, nil
}

func findConfigPath(appName string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	pkgDir := filepath.Dir(filename)
	projectRoot := filepath.Dir(filepath.Dir(pkgDir))

	return filepath.Join(projectRoot, "apps", appName, "config", "config.yaml")
}

func (i *baseAppInitializer) initLog() error {
	logCfg, ok := i.Log["default"]
	if !ok {
		for _, c := range i.Log {
			logCfg = c
			break
		}
	}
	return glog.InitLogger(&logCfg)
}

func (i *baseAppInitializer) initResources() error {
	var gormLogConfig *glog.LogConfig
	if c, ok := i.Log["gorm"]; ok {
		gormLogConfig = &c
	}
	if err := dbclient.InitMultiDB(i.DBConfigs, gormLogConfig); err != nil {
		return fmt.Errorf("init db: %w", err)
	}

	if i.RedisConfig.Addr != "" {
		var redisLogConfig *glog.LogConfig
		if c, ok := i.Log["redis"]; ok {
			redisLogConfig = &c
		}
		if err := dbclient.InitRedis(i.RedisConfig, redisLogConfig); err != nil {
			return fmt.Errorf("init redis: %w", err)
		}
	}

	var esLogConfig *glog.LogConfig
	if c, ok := i.Log["es"]; ok {
		esLogConfig = &c
	}
	if err := dbclient.InitMultiEs(i.ESConfigs, esLogConfig); err != nil {
		return fmt.Errorf("init elasticsearch: %w", err)
	}

	return nil
}

func (i *baseAppInitializer) Close() error {
	glog.Close()
	return nil
}
