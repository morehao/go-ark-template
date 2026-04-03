package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type RoleMenuCond struct {
	*genericdao.BaseCond
	RoleID   uint
	RoleIDs  []uint
	MenuID   uint
	TenantID uint
}

func (c *RoleMenuCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.RoleID > 0 {
		db.Where("role_id = ?", c.RoleID)
	}
	if len(c.RoleIDs) > 0 {
		db.Where("role_id IN (?)", c.RoleIDs)
	}
	if c.MenuID > 0 {
		db.Where("menu_id = ?", c.MenuID)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
}

type RoleMenuDao struct {
	*genericdao.GenericDao[iammodel.RoleMenuEntity, iammodel.RoleMenuEntityList]
}

func NewRoleMenuDao() *RoleMenuDao {
	return &RoleMenuDao{
		GenericDao: genericdao.NewGenericDao[iammodel.RoleMenuEntity, iammodel.RoleMenuEntityList](
			iammodel.TableNameRoleMenu, "RoleMenuDao",
			dbclient.IamDB,
		),
	}
}
