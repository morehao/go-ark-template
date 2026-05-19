package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type KnowledgeBaseCond struct {
	*genericdao.BaseCond
	TenantID uint
	Name     string
}

func (c *KnowledgeBaseCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.Name != "" {
		db.Where(tableName+".name = ?", c.Name)
	}
}

type KnowledgeBaseDao struct {
	*genericdao.GenericDao[model.KnowledgeBaseEntity, model.KnowledgeBaseEntityList]
}

func NewKnowledgeBaseDao() *KnowledgeBaseDao {
	return &KnowledgeBaseDao{
		GenericDao: genericdao.NewGenericDao[model.KnowledgeBaseEntity, model.KnowledgeBaseEntityList](
			model.TableNameKnowledgeBase, "KnowledgeBaseDao",
			dbclient.RagForgeDB,
		),
	}
}
