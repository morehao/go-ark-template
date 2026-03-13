package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type RoleCond struct {
	*genericdao.BaseCond
	CompanyID uint
	RoleName  string
	RoleCode  string
	Status    string
}

func (c *RoleCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.CompanyID > 0 {
		db.Where("company_id = ?", c.CompanyID)
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
	*genericdao.GenericDao[iammodel.RoleEntity, iammodel.RoleEntityList]
}

func NewRoleDao() *RoleDao {
	return &RoleDao{
		GenericDao: genericdao.NewGenericDao[iammodel.RoleEntity, iammodel.RoleEntityList](
			iammodel.TableNameRole, "RoleDao",
			dbclient.IamDB,
		),
	}
}
