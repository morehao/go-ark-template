package dao

import (
	"github.com/morehao/goark/apps/demo/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserCond struct {
	*genericdao.BaseCond
	CompanyID    uint
	CreatedBy    uint
	DeletedBy    uint
	DepartmentID uint
	Name         string
	UpdatedBy    uint
}

func (c *UserCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.CompanyID > 0 {
		db.Where("company_id = ?", c.CompanyID)
	}
	if c.CreatedBy > 0 {
		db.Where("created_by = ?", c.CreatedBy)
	}
	if c.DeletedBy > 0 {
		db.Where("deleted_by = ?", c.DeletedBy)
	}
	if c.DepartmentID > 0 {
		db.Where("department_id = ?", c.DepartmentID)
	}
	if c.Name != "" {
		db.Where("name = ?", c.Name)
	}
	if c.UpdatedBy > 0 {
		db.Where("updated_by = ?", c.UpdatedBy)
	}
}

type UserDao struct {
	*genericdao.GenericDao[model.UserEntity, model.UserEntityList]
}

func NewUserDao() *UserDao {
	return &UserDao{
		GenericDao: genericdao.NewGenericDao[model.UserEntity, model.UserEntityList](
			model.TableNameUser, "UserDao",
			dbclient.DemoDB,
		),
	}
}

// func (d *UserDao) WithTx(db *gorm.DB) *UserDao {
// 	return &UserDao{
// 		GenericDao: d.GenericDao.WithTx(db),
// 	}
// }
