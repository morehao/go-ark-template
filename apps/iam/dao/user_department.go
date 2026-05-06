package dao

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserDepartmentCond struct {
	*genericdao.BaseCond
	UserID   uint
	DeptID   uint
	TenantID uint
	DeptType model.UserDeptType
}

func (c *UserDepartmentCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.UserID > 0 {
		db.Where(tableName+".user_id = ?", c.UserID)
	}
	if c.DeptID > 0 {
		db.Where(tableName+".dept_id = ?", c.DeptID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.DeptType != "" {
		db.Where(tableName+".dept_type = ?", c.DeptType)
	}
}

type UserDepartmentDao struct {
	*genericdao.GenericDao[model.UserDepartmentEntity, model.UserDepartmentEntityList]
}

func NewUserDepartmentDao() *UserDepartmentDao {
	return &UserDepartmentDao{
		GenericDao: genericdao.NewGenericDao[model.UserDepartmentEntity, model.UserDepartmentEntityList](
			model.TableNameUserDepartment, "UserDepartmentDao",
			dbclient.IamDB,
		),
	}
}
