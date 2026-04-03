package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserDepartmentCond struct {
	*genericdao.BaseCond
	UserID   uint
	DeptID   uint
	TenantID uint
	DeptType iammodel.UserDeptType
}

func (c *UserDepartmentCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.UserID > 0 {
		db.Where("user_id = ?", c.UserID)
	}
	if c.DeptID > 0 {
		db.Where("dept_id = ?", c.DeptID)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.DeptType != "" {
		db.Where("dept_type = ?", c.DeptType)
	}
}

type UserDepartmentDao struct {
	*genericdao.GenericDao[iammodel.UserDepartmentEntity, iammodel.UserDepartmentEntityList]
}

func NewUserDepartmentDao() *UserDepartmentDao {
	return &UserDepartmentDao{
		GenericDao: genericdao.NewGenericDao[iammodel.UserDepartmentEntity, iammodel.UserDepartmentEntityList](
			iammodel.TableNameUserDepartment, "UserDepartmentDao",
			dbclient.IamDB,
		),
	}
}
