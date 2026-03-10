package dbclient

import (
	"fmt"

	"github.com/morehao/golib/dbaccess/dbgorm"
	"github.com/morehao/golib/glog"
	"gorm.io/gorm"
)

var (
	DBDemo *gorm.DB
	DBIam  *gorm.DB
)

const (
	DBNameDemo = "demoapp"
	DBNameIam  = "iam"
)

func InitMultiDB(configs []dbgorm.GormConfig, logConfig *glog.LogConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("mysql config is empty")
	}

	var opts []dbgorm.Option
	if logConfig != nil {
		opts = append(opts, dbgorm.WithLogConfig(logConfig))
	}
	for _, cfg := range configs {
		client, err := dbgorm.New(&cfg, opts...)
		if err != nil {
			return fmt.Errorf("init mysql failed: " + err.Error())
		}
		switch cfg.Service {
		case DBNameDemo:
			DBDemo = client
		case DBNameIam:
			DBIam = client
		default:
			return fmt.Errorf("unknown database service: " + cfg.Service)
		}
	}
	return nil
}
