package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type RoleMenuCond struct {
	*genericdao.BaseCond
	RoleID   uint
	MenuID   uint
	TenantID uint
}

func (c *RoleMenuCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.RoleID > 0 {
		db.Where(tableName+".role_id = ?", c.RoleID)
	}
	if c.MenuID > 0 {
		db.Where(tableName+".menu_id = ?", c.MenuID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
}

type RoleMenuDao struct {
	*genericdao.GenericDao[model.RoleMenuEntity, model.RoleMenuEntityList]
}

func NewRoleMenuDao() *RoleMenuDao {
	return &RoleMenuDao{
		GenericDao: genericdao.NewGenericDao[model.RoleMenuEntity, model.RoleMenuEntityList](
			model.TableNameRoleMenu, "RoleMenuDao",
			dbclient.IamDB,
		),
	}
}
