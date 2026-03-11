package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/pkg/genericdao"
	"gorm.io/gorm"
)

type RoleMenuCond struct {
	*genericdao.BaseCond
	RoleID    uint
	MenuID    uint
	CompanyID uint
}

func (c *RoleMenuCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.RoleID > 0 {
		db.Where("role_id = ?", c.RoleID)
	}
	if c.MenuID > 0 {
		db.Where("menu_id = ?", c.MenuID)
	}
	if c.CompanyID > 0 {
		db.Where("company_id = ?", c.CompanyID)
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

func (d *RoleMenuDao) WithTx(db *gorm.DB) *RoleMenuDao {
	return &RoleMenuDao{
		GenericDao: d.GenericDao.WithTx(db),
	}
}
