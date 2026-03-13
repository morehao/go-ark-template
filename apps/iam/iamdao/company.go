package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type CompanyCond struct {
	*genericdao.BaseCond
	TenantID    uint
	CompanyName string
	CompanyCode string
	Status      string
}

func (c *CompanyCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.CompanyName != "" {
		db.Where("company_name = ?", c.CompanyName)
	}
	if c.CompanyCode != "" {
		db.Where("company_code = ?", c.CompanyCode)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type CompanyDao struct {
	*genericdao.GenericDao[iammodel.CompanyEntity, iammodel.CompanyEntityList]
}

func NewCompanyDao() *CompanyDao {
	return &CompanyDao{
		GenericDao: genericdao.NewGenericDao[iammodel.CompanyEntity, iammodel.CompanyEntityList](
			iammodel.TableNameCompany, "CompanyDao",
			dbclient.IamDB,
		),
	}
}
