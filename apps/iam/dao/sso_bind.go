package dao

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type SSOBindCond struct {
	*genericdao.BaseCond
	OrgID    uint
	TenantID uint
	UserID   uint
	SSOType  string
	OpenID   string
}

func (c *SSOBindCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.OrgID > 0 {
		db.Where(tableName + ".org_id = ?", c.OrgID)
	}
	if c.TenantID > 0 {
		db.Where(tableName + ".tenant_id = ?", c.TenantID)
	}
	if c.UserID > 0 {
		db.Where(tableName + ".user_id = ?", c.UserID)
	}
	if c.SSOType != "" {
		db.Where(tableName + ".sso_type = ?", c.SSOType)
	}
	if c.OpenID != "" {
		db.Where(tableName + ".open_id = ?", c.OpenID)
	}
}

type SSOBindDao struct {
	*genericdao.GenericDao[model.SSOBindEntity, model.SSOBindEntityList]
}

func NewSSOBindDao() *SSOBindDao {
	return &SSOBindDao{
		GenericDao: genericdao.NewGenericDao[model.SSOBindEntity, model.SSOBindEntityList](
			model.TableNameSSOBind, "SSOBindDao",
			dbclient.IamDB,
		),
	}
}