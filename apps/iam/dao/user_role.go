package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserRoleCond struct {
	*genericdao.BaseCond
	UserID   uint
	RoleID   uint
	TenantID uint
}

func (c *UserRoleCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.UserID > 0 {
		db.Where("user_id = ?", c.UserID)
	}
	if c.RoleID > 0 {
		db.Where("role_id = ?", c.RoleID)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
}

type UserRoleDao struct {
	*genericdao.GenericDao[model.UserRoleEntity, model.UserRoleEntityList]
}

func NewUserRoleDao() *UserRoleDao {
	return &UserRoleDao{
		GenericDao: genericdao.NewGenericDao[model.UserRoleEntity, model.UserRoleEntityList](
			model.TableNameUserRole, "UserRoleDao",
			dbclient.IamDB,
		),
	}
}
