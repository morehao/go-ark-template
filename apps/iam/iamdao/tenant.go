package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type TenantCond struct {
	*genericdao.BaseCond
	Name string
}

func (c *TenantCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Name != "" {
		db.Where("name = ?", c.Name)
	}
}

type TenantDao struct {
	*genericdao.GenericDao[iammodel.TenantEntity, iammodel.TenantEntityList]
}

func NewTenantDao() *TenantDao {
	return &TenantDao{
		GenericDao: genericdao.NewGenericDao[iammodel.TenantEntity, iammodel.TenantEntityList](
			iammodel.TableNameTenant, "TenantDao",
			dbclient.IamDB,
		),
	}
}
