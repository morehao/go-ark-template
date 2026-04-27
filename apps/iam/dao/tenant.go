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
	ParentID   uint
	TenantName string
	TenantCode string
	Status     model.TenantStatus
	Domain     string
}

func (c *TenantCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrgID > 0 {
		db.Where(tableName+".org_id = ?", c.OrgID)
	}
	if c.ParentID > 0 {
		db.Where(tableName+".parent_id = ?", c.ParentID)
	}
	if c.TenantName != "" {
		db.Where(tableName+".tenant_name = ?", c.TenantName)
	}
	if c.TenantCode != "" {
		db.Where(tableName+".tenant_code = ?", c.TenantCode)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
	if c.Domain != "" {
		db.Where(tableName+".domain = ?", c.Domain)
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
