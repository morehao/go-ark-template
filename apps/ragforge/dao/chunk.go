package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type ChunkCond struct {
	*genericdao.BaseCond
	KnowledgeID uint
	KbID        uint
	TenantID    uint
}

func (c *ChunkCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.KnowledgeID > 0 {
		db.Where(tableName+".knowledge_id = ?", c.KnowledgeID)
	}
	if c.KbID > 0 {
		db.Where(tableName+".kb_id = ?", c.KbID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
}

type ChunkDao struct {
	*genericdao.GenericDao[model.ChunkEntity, model.ChunkEntityList]
}

func NewChunkDao() *ChunkDao {
	return &ChunkDao{
		GenericDao: genericdao.NewGenericDao[model.ChunkEntity, model.ChunkEntityList](
			model.TableNameChunk, "ChunkDao",
			dbclient.RagForgeDB,
		),
	}
}
