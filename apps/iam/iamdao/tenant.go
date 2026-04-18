package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type TenantCond struct {
	*genericdao.BaseCond
	OrganizationID uint
	TenantName     string
	TenantCode     string
	Status         iammodel.TenantStatus
}

func (c *TenantCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrganizationID > 0 {
		db.Where("organization_id = ?", c.OrganizationID)
	}
	if c.TenantName != "" {
		db.Where("tenant_name = ?", c.TenantName)
	}
	if c.TenantCode != "" {
		db.Where("tenant_code = ?", c.TenantCode)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type TenantDao struct {
	*genericdao.GenericDao[iammodel.TenantEntity, iammodel.TenantEntityList]
}

func NewTenantDao() *TenantDao {
	return &TenantDao{
		GenericDao: genericdao.NewGenericDao[iammodel.TenantEntity, iammodel.TenantEntityList](
			iammodel.TableNameTenant, "TenantDao",
			dbclient.IamDB,
		),
	}
}
