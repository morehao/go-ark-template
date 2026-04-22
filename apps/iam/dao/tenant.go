package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type TenantCond struct {
	*genericdao.BaseCond
	OrgID      uint
	TenantName string
	TenantCode string
	Status     model.TenantStatus
}

func (c *TenantCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrgID > 0 {
		db.Where("org_id = ?", c.OrgID)
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
	*genericdao.GenericDao[model.TenantEntity, model.TenantEntityList]
}

func NewTenantDao() *TenantDao {
	return &TenantDao{
		GenericDao: genericdao.NewGenericDao[model.TenantEntity, model.TenantEntityList](
			model.TableNameTenant, "TenantDao",
			dbclient.IamDB,
		),
	}
}