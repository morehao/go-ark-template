package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type VectorStoreCond struct {
	*genericdao.BaseCond
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
	*genericdao.GenericDao[model.VectorStoreEntity, model.VectorStoreEntityList]
}

func NewVectorStoreDao() *VectorStoreDao {
	return &VectorStoreDao{
		GenericDao: genericdao.NewGenericDao[model.VectorStoreEntity, model.VectorStoreEntityList](
			model.TableNameVectorStore, "VectorStoreDao",
			dbclient.RagForgeDB,
		),
	}
}
