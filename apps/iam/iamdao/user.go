package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/pkg/genericdao"
	"gorm.io/gorm"
)

type UserCond struct {
	*genericdao.BaseCond
	Username  string
	CompanyID uint
	Status    string
}

func (c *UserCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Username != "" {
		db.Where("username = ?", c.Username)
	}
	if c.CompanyID > 0 {
		db.Where("company_id = ?", c.CompanyID)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type UserDao struct {
	*genericdao.GenericDao[iammodel.UserEntity, iammodel.UserEntityList]
}

func NewUserDao() *UserDao {
	return &UserDao{
		GenericDao: genericdao.NewGenericDao[iammodel.UserEntity, iammodel.UserEntityList](
			iammodel.TableNameUser, "UserDao",
			dbclient.IamDB,
		),
	}
}

func (d *UserDao) WithTx(db *gorm.DB) *UserDao {
	return &UserDao{
		GenericDao: d.GenericDao.WithTx(db),
	}
}
