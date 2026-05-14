package dao

import (
	"github.com/morehao/goark/apps/ragflow/model"
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type KnowledgeBaseCond struct {
	*genericdao.BaseCond
	Name     string
	Status   string
}

func (c *KnowledgeBaseCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.Name != "" {
		db.Where("name = ?", c.Name)
	}
	if c.Status != "" {
		db.Where("status = ?", c.Status)
	}
}

type KnowledgeBaseDao struct {
	*genericdao.GenericDao[model.KnowledgeBaseEntity, model.KnowledgeBaseEntityList]
}

func NewKnowledgeBaseDao() *KnowledgeBaseDao {
	return &KnowledgeBaseDao{
		GenericDao: genericdao.NewGenericDao[model.KnowledgeBaseEntity, model.KnowledgeBaseEntityList](
			model.TableNameKnowledgeBase, "KnowledgeBaseDao",
			dbclient.RagflowDB,
		),
	}
}