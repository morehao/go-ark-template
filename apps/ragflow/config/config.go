package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/morehao/golib/conf"
	"github.com/morehao/golib/dbaccess/dbes"
	"github.com/morehao/golib/dbaccess/dbgorm"
	"github.com/morehao/golib/dbaccess/dbredis"
	"github.com/morehao/golib/glog"
	"github.com/morehao/golib/gtrace"
	"github.com/morehao/golib/protocol/ghttp"
)

var Conf *Config

type Config struct {
	Server      Server                    `yaml:"server"`
	Log         map[string]glog.LogConfig `yaml:"log"`
	Trace       gtrace.TraceConfig        `yaml:"trace"`
	DBConfigs   []dbgorm.GormConfig       `yaml:"db_configs"`
	RedisConfig dbredis.RedisConfig       `yaml:"redis_config"`
	ESConfigs   []dbes.ESConfig           `yaml:"es_configs"`
	Client      Client                    `yaml:"client"`
}

type Server struct {
	Name string `yaml:"name"` // 服务名称
	Port string `yaml:"port"` // 服务端口
	Env  string `yaml:"env"`  // 环境变量
}

type Client struct {
	HTTPBingo *ghttp.Client `yaml:"httpbingo"`
}

func InitConf() {
	configPath := getConfigPath()
	LoadConfig(configPath)
}

func LoadConfig(configPath string) {
	fmt.Println("Load config file:", configPath)

	var cfg Config
	conf.LoadConfig(configPath, &cfg)
	Conf = &cfg
}

func getConfigPath() string {
	if configPath := os.Getenv("APP_CONFIG_PATH"); configPath != "" {
		return configPath
	}

	relativePath := "../config/config.yaml"
	if fileExists(relativePath) {
		return relativePath
	}

	execPath, err := os.Executable()
	if err == nil {
		absPath := filepath.Join(filepath.Dir(execPath), "..", "config", "config.yaml")
		if fileExists(absPath) {
			return absPath
		}
	}

	return relativePath
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}