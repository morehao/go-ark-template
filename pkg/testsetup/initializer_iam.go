package testsetup

import (
	"fmt"
	"os"

	"github.com/morehao/goark/apps/iam/config"
	"github.com/morehao/golib/biz/testkit"
)

type iamappInitializer struct {
	*baseAppInitializer
	*testkit.BaseInitializer
}

func newIamappInitializer() (Initializer, error) {
	base, err := newBaseAppInitializer(AppNameIam)
	if err != nil {
		return nil, err
	}

	baseInit, err := testkit.NewBaseInitializer(AppNameIam)
	if err != nil {
		return nil, err
	}

	return &iamappInitializer{
		baseAppInitializer: base,
		BaseInitializer:    baseInit,
	}, nil
}

func (i *iamappInitializer) Initialize() error {
	if _, err := os.Stat(i.ConfigPath); err != nil {
		return fmt.Errorf("config file not found: %s, error: %w", i.ConfigPath, err)
	}

	var panicErr interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
			}
		}()
		config.LoadConfig(i.ConfigPath)
	}()

	if panicErr != nil {
		return fmt.Errorf("load config failed: %v", panicErr)
	}

	i.Log = config.Conf.Log
	i.DBConfigs = config.Conf.DBConfigs
	i.RedisConfig = config.Conf.RedisConfig
	i.ESConfigs = config.Conf.ESConfigs

	if err := i.initLog(); err != nil {
		return fmt.Errorf("init log: %w", err)
	}

	if err := i.initResources(); err != nil {
		return fmt.Errorf("init resources: %w", err)
	}

	return nil
}
