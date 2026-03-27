package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrganizationConfigCond struct {
	*genericdao.BaseCond
	OrganizationID uint
	ConfigGroup    string
	ConfigKey      string
}

func (c *OrganizationConfigCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrganizationID > 0 {
		db.Where("organization_id = ?", c.OrganizationID)
	}
	if c.ConfigGroup != "" {
		db.Where("config_group = ?", c.ConfigGroup)
	}
	if c.ConfigKey != "" {
		db.Where("config_key = ?", c.ConfigKey)
	}
}

type OrganizationConfigDao struct {
	*genericdao.GenericDao[iammodel.OrganizationConfigEntity, iammodel.OrganizationConfigEntityList]
}

func NewOrganizationConfigDao() *OrganizationConfigDao {
	return &OrganizationConfigDao{
		GenericDao: genericdao.NewGenericDao[iammodel.OrganizationConfigEntity, iammodel.OrganizationConfigEntityList](
			iammodel.TableNameOrganizationConfig, "OrganizationConfigDao",
			dbclient.IamDB,
		),
	}
}
