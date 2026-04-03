package testsetup

import (
	"fmt"
	"os"

	"github.com/morehao/goark/apps/demo/config"
	"github.com/morehao/goark/pkg/testkit"
)

type demoappInitializer struct {
	*baseAppInitializer
	*testkit.BaseInitializer
}

func newDemoappInitializer() (Initializer, error) {
	base, err := newBaseAppInitializer(AppNameDemo)
	if err != nil {
		return nil, err
	}

	baseInit, err := testkit.NewBaseInitializer(AppNameDemo)
	if err != nil {
		return nil, err
	}

	return &demoappInitializer{
		baseAppInitializer: base,
		BaseInitializer:    baseInit,
	}, nil
}

func (d *demoappInitializer) Initialize() error {
	if _, err := os.Stat(d.ConfigPath); err != nil {
		return fmt.Errorf("config file not found: %s, error: %w", d.ConfigPath, err)
	}

	var panicErr interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
			}
		}()
		config.LoadConfig(d.ConfigPath)
	}()

	if panicErr != nil {
		return fmt.Errorf("load config failed: %v", panicErr)
	}

	d.Log = config.Conf.Log
	d.DBConfigs = config.Conf.DBConfigs
	d.RedisConfig = config.Conf.RedisConfig
	d.ESConfigs = config.Conf.ESConfigs

	if err := d.initLog(); err != nil {
		return fmt.Errorf("init log: %w", err)
	}

	if err := d.initResources(); err != nil {
		return fmt.Errorf("init resources: %w", err)
	}

	return nil
}
