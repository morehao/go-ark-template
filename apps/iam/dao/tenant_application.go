package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type TenantApplicationCond struct {
	*genericdao.BaseCond
	AppID     uint
	CreatedBy uint
	TenantID  uint
}

func (c *TenantApplicationCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.AppID != 0 {
		db.Where(tableName+".app_id = ?", c.AppID)
	}
	if c.CreatedBy != 0 {
		db.Where(tableName+".created_by = ?", c.CreatedBy)
	}
	if c.TenantID != 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
}

type TenantApplicationDao struct {
	*genericdao.GenericDao[model.TenantApplicationEntity, model.TenantApplicationEntityList]
}

func NewTenantApplicationDao() *TenantApplicationDao {
	return &TenantApplicationDao{
		GenericDao: genericdao.NewGenericDao[model.TenantApplicationEntity, model.TenantApplicationEntityList](
			model.TableNameTenantApplication, "TenantApplicationDao",
			dbclient.IamDB,
		),
	}
}
