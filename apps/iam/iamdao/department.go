package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/pkg/genericdao"
	"gorm.io/gorm"
)

type DepartmentCond struct {
	*genericdao.BaseCond
	CompanyID uint
	ParentID  uint
	DeptName  string
	DeptCode  string
	Status    string
}

func (c *DepartmentCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.CompanyID > 0 {
		db.Where("company_id = ?", c.CompanyID)
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

func (d *DepartmentDao) WithTx(db *gorm.DB) *DepartmentDao {
	return &DepartmentDao{
		GenericDao: d.GenericDao.WithTx(db),
	}
}
