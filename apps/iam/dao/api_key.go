package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type ApiKeyCond struct {
	*genericdao.BaseCond
	ID         uint
	TenantID   uint
	UserID     uint
	AppID      uint
	KeyPrefix  string
	KeyName    string
	Status     model.ApiKeyStatus
}

func (c *ApiKeyCond) BuildCondition(db *gorm.DB, tableName string) {
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
	if c.AppID > 0 {
		db.Where(tableName+".app_id = ?", c.AppID)
	}
	if c.KeyPrefix != "" {
		db.Where(tableName+".key_prefix = ?", c.KeyPrefix)
	}
	if c.KeyName != "" {
		db.Where(tableName+".key_name = ?", c.KeyName)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type ApiKeyDao struct {
	*genericdao.GenericDao[model.ApiKeyEntity, model.ApiKeyEntityList]
}

func NewApiKeyDao() *ApiKeyDao {
	return &ApiKeyDao{
		GenericDao: genericdao.NewGenericDao[model.ApiKeyEntity, model.ApiKeyEntityList](
			model.TableNameApiKey, "ApiKeyDao",
			dbclient.IamDB,
		),
	}
}
