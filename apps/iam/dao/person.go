package dao

import (
	"github.com/morehao/goark/iam/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type PersonCond struct {
	*genericdao.BaseCond
	Mobile   string
	Email    string
	RealName string
}

func (c *PersonCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Mobile != "" {
		db.Where(tableName+".mobile = ?", c.Mobile)
	}
	if c.Email != "" {
		db.Where(tableName+".email = ?", c.Email)
	}
	if c.RealName != "" {
		db.Where(tableName+".real_name = ?", c.RealName)
	}
}

type PersonDao struct {
	*genericdao.GenericDao[model.PersonEntity, model.PersonEntityList]
}

func NewPersonDao() *PersonDao {
	return &PersonDao{
		GenericDao: genericdao.NewGenericDao[model.PersonEntity, model.PersonEntityList](
			model.TableNamePerson, "PersonDao",
			dbclient.IamDB,
		),
	}
}
