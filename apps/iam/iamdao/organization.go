package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OrganizationCond struct {
	*genericdao.BaseCond
	ID     uint
	Name   string
	Domain string
	Status iammodel.OrgStatus
}

func (c *OrganizationCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Name != "" {
		db.Where("organization_name = ?", c.Name)
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

type OrganizationDao struct {
	*genericdao.GenericDao[iammodel.OrganizationEntity, iammodel.OrganizationEntityList]
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{
		GenericDao: genericdao.NewGenericDao[iammodel.OrganizationEntity, iammodel.OrganizationEntityList](
			iammodel.TableNameOrganization, "OrganizationDao",
			dbclient.IamDB,
		),
	}
}
