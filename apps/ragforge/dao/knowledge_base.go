package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type KnowledgeBaseCond struct {
	*gormdao.BaseCond
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
	*gormdao.Dao[model.KnowledgeBaseEntity, model.KnowledgeBaseEntityList]
}

func NewKnowledgeBaseDao() *KnowledgeBaseDao {
	return &KnowledgeBaseDao{
		Dao: gormdao.NewDao[model.KnowledgeBaseEntity, model.KnowledgeBaseEntityList](
			model.TableNameKnowledgeBase, "KnowledgeBaseDao",
			dbclient.RagForgeDB,
		),
	}
}
