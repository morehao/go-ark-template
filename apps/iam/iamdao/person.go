package iamdao

import (
	"github.com/morehao/goark/apps/iam/iammodel"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type PersonCond struct {
	*genericdao.BaseCond
	Mobile   string
	RealName string
}

func (c *PersonCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Mobile != "" {
		db.Where("mobile = ?", c.Mobile)
	}
	if c.RealName != "" {
		db.Where("real_name = ?", c.RealName)
	}
}

type PersonDao struct {
	*genericdao.GenericDao[iammodel.PersonEntity, iammodel.PersonEntityList]
}

func NewPersonDao() *PersonDao {
	return &PersonDao{
		GenericDao: genericdao.NewGenericDao[iammodel.PersonEntity, iammodel.PersonEntityList](
			iammodel.TableNamePerson, "PersonDao",
			dbclient.IamDB,
		),
	}
}

func (d *PersonDao) WithTx(db *gorm.DB) *PersonDao {
	return &PersonDao{
		GenericDao: d.GenericDao.WithTx(db),
	}
}
