package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrgCond struct {
	*genericdao.BaseCond
	ID     uint
	Name   string
	Domain string
	Status model.OrgStatus
}

func (c *OrgCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Name != "" {
		db.Where("org_name = ?", c.Name)
	}
	if c.ID > 0 {
		db.Where("id = ?", c.ID)
	}
	if c.Domain != "" {
		db.Where("domain = ?", c.Domain)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type OrgDao struct {
	*genericdao.GenericDao[model.OrgEntity, model.OrgEntityList]
}

func NewOrgDao() *OrgDao {
	return &OrgDao{
		GenericDao: genericdao.NewGenericDao[model.OrgEntity, model.OrgEntityList](
			model.TableNameOrg, "OrgDao",
			dbclient.IamDB,
		),
	}
}