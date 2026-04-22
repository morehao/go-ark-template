package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrganizationApplicationCond struct {
	*genericdao.BaseCond
	AppID     uint
	CreatedBy uint
	OrgID     uint
}

func (c *OrganizationApplicationCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.AppID != 0 {
		db.Where(tableName+".app_id = ?", c.AppID)
	}
	if c.CreatedBy != 0 {
		db.Where(tableName+".created_by = ?", c.CreatedBy)
	}
	if c.OrgID != 0 {
		db.Where(tableName+".org_id = ?", c.OrgID)
	}
}

type OrganizationApplicationDao struct {
	*genericdao.GenericDao[model.OrganizationApplicationEntity, model.OrganizationApplicationEntityList]
}

func NewOrganizationApplicationDao() *OrganizationApplicationDao {
	return &OrganizationApplicationDao{
		GenericDao: genericdao.NewGenericDao[model.OrganizationApplicationEntity, model.OrganizationApplicationEntityList](
			model.TableNameOrganizationApplication, "OrganizationApplicationDao",
			dbclient.IamDB,
		),
	}
}
