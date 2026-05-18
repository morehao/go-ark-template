package dao

import (
	"context"

	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserCond struct {
	*genericdao.BaseCond
	Username string
	TenantID uint
	PersonID uint
	Status   model.UserStatus
}

func (c *UserCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Username != "" {
		db.Where(tableName+".username = ?", c.Username)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.PersonID > 0 {
		db.Where(tableName+".person_id = ?", c.PersonID)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type UserDao struct {
	*genericdao.GenericDao[model.UserEntity, model.UserEntityList]
}

func NewUserDao() *UserDao {
	return &UserDao{
		GenericDao: genericdao.NewGenericDao[model.UserEntity, model.UserEntityList](
			model.TableNameUser, "UserDao",
			dbclient.IamDB,
		),
	}
}

func (dao *UserDao) GetPendingUsers(ctx context.Context, tenantID uint) (model.UserEntityList, error) {
	result, err := dao.GetListByCond(ctx, &UserCond{
		TenantID: tenantID,
		Status:   model.UserStatusPending,
	})
	return result, err
}
