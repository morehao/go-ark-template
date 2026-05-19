package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type ModelCond struct {
	*genericdao.BaseCond
	TenantID  uint
	Name      string
	ModelType model.ModelType
	Provider  string
	Status    model.ModelStatus
}

func (c *ModelCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.Name != "" {
		db.Where(tableName+".name like ?", "%"+c.Name+"%")
	}
	if c.ModelType != "" {
		db.Where(tableName+".model_type = ?", c.ModelType)
	}
	if c.Provider != "" {
		db.Where(tableName+".provider = ?", c.Provider)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type ModelDao struct {
	*genericdao.GenericDao[model.ModelEntity, model.ModelEntityList]
}

func NewModelDao() *ModelDao {
	return &ModelDao{
		GenericDao: genericdao.NewGenericDao[model.ModelEntity, model.ModelEntityList](
			model.TableNameModel, "ModelDao",
			dbclient.RagForgeDB,
		),
	}
}
