package dao

import (
	"github.com/morehao/goark/apps/ragflow/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type DocumentCond struct {
	*genericdao.BaseCond
	KnowledgeBaseID uint
	Name            string
	Type            string
	Status          string
}

func (c *DocumentCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.KnowledgeBaseID > 0 {
		db.Where("knowledge_base_id = ?", c.KnowledgeBaseID)
	}
	if c.Name != "" {
		db.Where("name = ?", c.Name)
	}
	if c.Type != "" {
		db.Where("type = ?", c.Type)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type DocumentDao struct {
	*genericdao.GenericDao[model.DocumentEntity, model.DocumentEntityList]
}

func NewDocumentDao() *DocumentDao {
	return &DocumentDao{
		GenericDao: genericdao.NewGenericDao[model.DocumentEntity, model.DocumentEntityList](
			model.TableNameDocument, "DocumentDao",
			dbclient.RagflowDB,
		),
	}
}