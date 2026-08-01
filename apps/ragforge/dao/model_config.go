package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type ModelCond struct {
	*gormdao.BaseCond
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
	*gormdao.Dao[model.ModelEntity, model.ModelEntityList]
}

func NewModelDao() *ModelDao {
	return &ModelDao{
		Dao: gormdao.NewDao[model.ModelEntity, model.ModelEntityList](
			model.TableNameModel, "ModelDao",
			dbclient.RagForgeDB,
		),
	}
}
