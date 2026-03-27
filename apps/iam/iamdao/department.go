package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type DepartmentCond struct {
	*genericdao.BaseCond
	TenantID uint
	ParentID uint
	DeptName string
	DeptCode string
	Status   string
}

func (c *DepartmentCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.ParentID > 0 {
		db.Where("parent_id = ?", c.ParentID)
	}
	if c.DeptName != "" {
		db.Where("dept_name = ?", c.DeptName)
	}
	if c.DeptCode != "" {
		db.Where("dept_code = ?", c.DeptCode)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type DepartmentDao struct {
	*genericdao.GenericDao[iammodel.DepartmentEntity, iammodel.DepartmentEntityList]
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{
		GenericDao: genericdao.NewGenericDao[iammodel.DepartmentEntity, iammodel.DepartmentEntityList](
			iammodel.TableNameDepartment, "DepartmentDao",
			dbclient.IamDB,
		),
	}
}
