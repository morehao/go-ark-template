package dao

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type MenuCond struct {
	*genericdao.BaseCond
	TenantID       uint
	ParentID       uint
	MenuName       string
	MenuCode       string
	MenuType       model.MenuType
	Status         model.MenuStatus
	AccessPolicies []model.MenuAccessPolicyString
}

func (c *MenuCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.ParentID > 0 {
		db.Where(tableName+".parent_id = ?", c.ParentID)
	}
	if c.MenuName != "" {
		db.Where(tableName+".menu_name = ?", c.MenuName)
	}
	if c.MenuCode != "" {
		db.Where(tableName+".menu_code = ?", c.MenuCode)
	}
	if c.MenuType != "" {
		db.Where(tableName+".menu_type = ?", c.MenuType)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
	if len(c.AccessPolicies) > 0 {
		mask := model.AccessPoliciesToMask(c.AccessPolicies)
		db.Where(tableName+".access_policy & ? != 0", mask)
	}
}

type MenuDao struct {
	*genericdao.GenericDao[model.MenuEntity, model.MenuEntityList]
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		GenericDao: genericdao.NewGenericDao[model.MenuEntity, model.MenuEntityList](
			model.TableNameMenu, "MenuDao",
			dbclient.IamDB,
		),
	}
}
