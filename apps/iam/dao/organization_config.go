package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrganizationConfigCond struct {
	*genericdao.BaseCond
	OrgID       uint
	ConfigGroup string
	ConfigKey   string
}

func (c *OrganizationConfigCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrgID > 0 {
		db.Where(tableName+".org_id = ?", c.OrgID)
	}
	if c.ConfigGroup != "" {
		db.Where(tableName+".config_group = ?", c.ConfigGroup)
	}
	if c.ConfigKey != "" {
		db.Where(tableName+".config_key = ?", c.ConfigKey)
	}
}

type OrganizationConfigDao struct {
	*genericdao.GenericDao[model.OrganizationConfigEntity, model.OrganizationConfigEntityList]
}

func NewOrganizationConfigDao() *OrganizationConfigDao {
	return &OrganizationConfigDao{
		GenericDao: genericdao.NewGenericDao[model.OrganizationConfigEntity, model.OrganizationConfigEntityList](
			model.TableNameOrganizationConfig, "OrganizationConfigDao",
			dbclient.IamDB,
		),
	}
}
