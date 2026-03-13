package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type UserRoleCond struct {
	*genericdao.BaseCond
	UserID    uint
	RoleID    uint
	CompanyID uint
}

func (c *UserRoleCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.UserID > 0 {
		db.Where("user_id = ?", c.UserID)
	}
	if c.RoleID > 0 {
		db.Where("role_id = ?", c.RoleID)
	}
	if c.CompanyID > 0 {
		db.Where("company_id = ?", c.CompanyID)
	}
}

type UserRoleDao struct {
	*genericdao.GenericDao[iammodel.UserRoleEntity, iammodel.UserRoleEntityList]
}

func NewUserRoleDao() *UserRoleDao {
	return &UserRoleDao{
		GenericDao: genericdao.NewGenericDao[iammodel.UserRoleEntity, iammodel.UserRoleEntityList](
			iammodel.TableNameUserRole, "UserRoleDao",
			dbclient.IamDB,
		),
	}
}
