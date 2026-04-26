package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type LoginLogCond struct {
	*genericdao.BaseCond
	ID          uint
	TenantID    uint
	UserID      uint
	Username    string
	LoginType   string
	LoginStatus string
}

func (c *LoginLogCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.ID > 0 {
		db.Where(tableName+".id = ?", c.ID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.UserID > 0 {
		db.Where(tableName+".user_id = ?", c.UserID)
	}
	if c.Username != "" {
		db.Where(tableName+".username = ?", c.Username)
	}
	if c.LoginType != "" {
		db.Where(tableName+".login_type = ?", c.LoginType)
	}
	if c.LoginStatus != "" {
		db.Where(tableName+".login_status = ?", c.LoginStatus)
	}
}

type LoginLogDao struct {
	*genericdao.GenericDao[model.LoginLogEntity, model.LoginLogEntityList]
}

func NewLoginLogDao() *LoginLogDao {
	return &LoginLogDao{
		GenericDao: genericdao.NewGenericDao[model.LoginLogEntity, model.LoginLogEntityList](
			model.TableNameLoginLog, "LoginLogDao",
			dbclient.IamDB,
		),
	}
}
