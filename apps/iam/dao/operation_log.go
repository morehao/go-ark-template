package dao

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type OperationLogCond struct {
	*genericdao.BaseCond
	ID         uint
	TenantID   uint
	UserID     uint
	Username   string
	Module     string
	Operation  string
	Status     string
	RequestID  string
}

func (c *OperationLogCond) BuildCondition(db *gorm.DB, tableName string) {
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
	if c.Module != "" {
		db.Where(tableName+".module = ?", c.Module)
	}
	if c.Operation != "" {
		db.Where(tableName+".operation = ?", c.Operation)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
	if c.RequestID != "" {
		db.Where(tableName+".request_id = ?", c.RequestID)
	}
}

type OperationLogDao struct {
	*genericdao.GenericDao[model.OperationLogEntity, model.OperationLogEntityList]
}

func NewOperationLogDao() *OperationLogDao {
	return &OperationLogDao{
		GenericDao: genericdao.NewGenericDao[model.OperationLogEntity, model.OperationLogEntityList](
			model.TableNameOperationLog, "OperationLogDao",
			dbclient.IamDB,
		),
	}
}
