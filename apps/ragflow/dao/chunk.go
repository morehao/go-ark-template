package dao

import (
	"github.com/morehao/goark/apps/ragflow/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type ChunkCond struct {
	*genericdao.BaseCond
	DocumentID uint
}

func (c *ChunkCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.DocumentID > 0 {
		db.Where("document_id = ?", c.DocumentID)
	}
}

type ChunkDao struct {
	*genericdao.GenericDao[model.ChunkEntity, model.ChunkEntityList]
}

func NewChunkDao() *ChunkDao {
	return &ChunkDao{
		GenericDao: genericdao.NewGenericDao[model.ChunkEntity, model.ChunkEntityList](
			model.TableNameChunk, "ChunkDao",
			dbclient.RagflowDB,
		),
	}
}