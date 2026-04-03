package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type MenuCond struct {
	*genericdao.BaseCond
	TenantID uint
	ParentID uint
	MenuName string
	MenuCode string
	MenuType string
	Status   string
}

func (c *MenuCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.ParentID > 0 {
		db.Where("parent_id = ?", c.ParentID)
	}
	if c.MenuName != "" {
		db.Where("menu_name = ?", c.MenuName)
	}
	if c.MenuCode != "" {
		db.Where("menu_code = ?", c.MenuCode)
	}
	if c.MenuType != "" {
		db.Where("menu_type = ?", c.MenuType)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type MenuDao struct {
	*genericdao.GenericDao[iammodel.MenuEntity, iammodel.MenuEntityList]
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		GenericDao: genericdao.NewGenericDao[iammodel.MenuEntity, iammodel.MenuEntityList](
			iammodel.TableNameMenu, "MenuDao",
			dbclient.IamDB,
		),
	}
}
