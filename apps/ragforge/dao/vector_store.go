package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type VectorStoreCond struct {
	*gormdao.BaseCond
	TenantID   uint
	Name       string
	EngineType model.EngineType
	Status     model.VectorStoreStatus
}

func (c *VectorStoreCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.Name != "" {
		db.Where(tableName+".name like ?", "%"+c.Name+"%")
	}
	if c.EngineType != "" {
		db.Where(tableName+".engine_type = ?", c.EngineType)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type VectorStoreDao struct {
	*gormdao.Dao[model.VectorStoreEntity, model.VectorStoreEntityList]
}

func NewVectorStoreDao() *VectorStoreDao {
	return &VectorStoreDao{
		Dao: gormdao.NewDao[model.VectorStoreEntity, model.VectorStoreEntityList](
			model.TableNameVectorStore, "VectorStoreDao",
			dbclient.RagForgeDB,
		),
	}
}
