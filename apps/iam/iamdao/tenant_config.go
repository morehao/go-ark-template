package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type TenantConfigCond struct {
	*genericdao.BaseCond
	TenantID    uint
	ConfigGroup string
	ConfigKey   string
}

func (c *TenantConfigCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.ConfigGroup != "" {
		db.Where("config_group = ?", c.ConfigGroup)
	}
	if c.ConfigKey != "" {
		db.Where("config_key = ?", c.ConfigKey)
	}
}

type TenantConfigDao struct {
	*genericdao.GenericDao[iammodel.TenantConfigEntity, iammodel.TenantConfigEntityList]
}

func NewTenantConfigDao() *TenantConfigDao {
	return &TenantConfigDao{
		GenericDao: genericdao.NewGenericDao[iammodel.TenantConfigEntity, iammodel.TenantConfigEntityList](
			iammodel.TableNameTenantConfig, "TenantConfigDao",
			dbclient.IamDB,
		),
	}
}

func (d *TenantConfigDao) WithTx(db *gorm.DB) *TenantConfigDao {
	return &TenantConfigDao{
		GenericDao: d.GenericDao.WithTx(db),
	}
}
