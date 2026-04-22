package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrganizationCond struct {
	*genericdao.BaseCond
	ID     uint
	Name   string
	Domain string
	Status model.OrgStatus
}

func (c *OrganizationCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Name != "" {
		db.Where(tableName+".org_name = ?", c.Name)
	}
	if c.ID > 0 {
		db.Where(tableName+".id = ?", c.ID)
	}
	if c.Domain != "" {
		db.Where(tableName+".domain = ?", c.Domain)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type OrganizationDao struct {
	*genericdao.GenericDao[model.OrganizationEntity, model.OrganizationEntityList]
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{
		GenericDao: genericdao.NewGenericDao[model.OrganizationEntity, model.OrganizationEntityList](
			model.TableNameOrganization, "OrganizationDao",
			dbclient.IamDB,
		),
	}
}
