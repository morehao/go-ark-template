package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type KnowledgeCond struct {
	*gormdao.BaseCond
	KbID          uint
	TenantID      uint
	KnowledgeType model.KnowledgeType
	ParseStatus   model.ParseStatus
}

func (c *KnowledgeCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.KbID > 0 {
		db.Where(tableName+".kb_id = ?", c.KbID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.KnowledgeType != "" {
		db.Where(tableName+".type = ?", c.KnowledgeType)
	}
	if c.ParseStatus != "" {
		db.Where(tableName+".parse_status = ?", c.ParseStatus)
	}
}

type KnowledgeDao struct {
	*gormdao.Dao[model.KnowledgeEntity, model.KnowledgeEntityList]
}

func NewKnowledgeDao() *KnowledgeDao {
	return &KnowledgeDao{
		Dao: gormdao.NewDao[model.KnowledgeEntity, model.KnowledgeEntityList](
			model.TableNameKnowledge, "KnowledgeDao",
			dbclient.RagForgeDB,
		),
	}
}
