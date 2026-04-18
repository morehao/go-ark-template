package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserCond struct {
	*genericdao.BaseCond
	Username string
	TenantID uint
	PersonID uint
	Status   iammodel.UserStatus
}

func (c *UserCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Username != "" {
		db.Where("username = ?", c.Username)
	}
	if c.TenantID > 0 {
		db.Where("tenant_id = ?", c.TenantID)
	}
	if c.PersonID > 0 {
		db.Where("person_id = ?", c.PersonID)
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
