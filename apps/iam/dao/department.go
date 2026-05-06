package dao

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type DepartmentCond struct {
	*genericdao.BaseCond
	TenantID    uint
	ParentID    uint
	ParentIDNil bool // 是否显式查询 ParentID（包括 ParentID=0 的情况）
	DeptName    string
	DeptCode    string
	Status      model.DeptStatus
}

func (c *DepartmentCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.ParentIDNil {
		db.Where(tableName+".parent_id = ?", c.ParentID)
	} else if c.ParentID > 0 {
		db.Where(tableName+".parent_id = ?", c.ParentID)
	}
	if c.DeptName != "" {
		db.Where(tableName+".dept_name = ?", c.DeptName)
	}
	if c.DeptCode != "" {
		db.Where(tableName+".dept_code = ?", c.DeptCode)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type DepartmentDao struct {
	*genericdao.GenericDao[model.DepartmentEntity, model.DepartmentEntityList]
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{
		GenericDao: genericdao.NewGenericDao[model.DepartmentEntity, model.DepartmentEntityList](
			model.TableNameDepartment, "DepartmentDao",
			dbclient.IamDB,
		),
	}
}
