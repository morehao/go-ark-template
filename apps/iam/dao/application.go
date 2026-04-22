package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type ApplicationCond struct {
	*genericdao.BaseCond
	AppCode     string
	AppName     string
	AppType     string
	CallbackUrl string
	CreatedBy   uint
	DeletedBy   uint
	Description string
	HomepageUrl string
	Logo        string
	SortOrder   int32
	Status      string
	UpdatedBy   uint
}

func (c *ApplicationCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.AppCode != "" {
		db.Where(tableName+".app_code = ?", c.AppCode)
	}
	if c.AppName != "" {
		db.Where(tableName+".app_name = ?", c.AppName)
	}
	if c.AppType != "" {
		db.Where(tableName+".app_type = ?", c.AppType)
	}
	if c.CallbackUrl != "" {
		db.Where(tableName+".callback_url = ?", c.CallbackUrl)
	}
	if c.CreatedBy != 0 {
		db.Where(tableName+".created_by = ?", c.CreatedBy)
	}
	if c.DeletedBy != 0 {
		db.Where(tableName+".deleted_by = ?", c.DeletedBy)
	}
	if c.Description != "" {
		db.Where(tableName+".description = ?", c.Description)
	}
	if c.HomepageUrl != "" {
		db.Where(tableName+".homepage_url = ?", c.HomepageUrl)
	}
	if c.Logo != "" {
		db.Where(tableName+".logo = ?", c.Logo)
	}
	if c.SortOrder != 0 {
		db.Where(tableName+".sort_order = ?", c.SortOrder)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
	if c.UpdatedBy != 0 {
		db.Where(tableName+".updated_by = ?", c.UpdatedBy)
	}
}

type ApplicationDao struct {
	*genericdao.GenericDao[model.ApplicationEntity, model.ApplicationEntityList]
}

func NewApplicationDao() *ApplicationDao {
	return &ApplicationDao{
		GenericDao: genericdao.NewGenericDao[model.ApplicationEntity, model.ApplicationEntityList](
			model.TableNameApplication, "ApplicationDao",
			dbclient.IamDB,
		),
	}
}
