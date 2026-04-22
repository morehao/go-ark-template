package dao

import (
	"github.com/morehao/goark/apps/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type RoleCond struct {
	*genericdao.BaseCond
	TenantID uint
	RoleName string
	RoleCode string
	Status   model.RoleStatus
}

func (c *RoleCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.RoleName != "" {
		db.Where("role_name = ?", c.RoleName)
	}
	if c.RoleCode != "" {
		db.Where("role_code = ?", c.RoleCode)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type RoleDao struct {
	*genericdao.GenericDao[model.RoleEntity, model.RoleEntityList]
}

func NewRoleDao() *RoleDao {
	return &RoleDao{
		GenericDao: genericdao.NewGenericDao[model.RoleEntity, model.RoleEntityList](
			model.TableNameRole, "RoleDao",
			dbclient.IamDB,
		),
	}
}
