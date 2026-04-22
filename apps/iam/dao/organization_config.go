package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrgConfigCond struct {
	*genericdao.BaseCond
	OrgID       uint
	ConfigGroup string
	ConfigKey   string
}

func (c *OrgConfigCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrgID > 0 {
		db.Where("org_id = ?", c.OrgID)
	}
	if c.ConfigGroup != "" {
		db.Where("config_group = ?", c.ConfigGroup)
	}
	if c.ConfigKey != "" {
		db.Where("config_key = ?", c.ConfigKey)
	}
}

type OrgConfigDao struct {
	*genericdao.GenericDao[model.OrgConfigEntity, model.OrgConfigEntityList]
}

func NewOrgConfigDao() *OrgConfigDao {
	return &OrgConfigDao{
		GenericDao: genericdao.NewGenericDao[model.OrgConfigEntity, model.OrgConfigEntityList](
			model.TableNameOrgConfig, "OrgConfigDao",
			dbclient.IamDB,
		),
	}
}