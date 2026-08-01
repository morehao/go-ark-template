package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type FAQCond struct {
	*gormdao.BaseCond
	KbID     uint
	TenantID uint
	Question string
	Status   model.FAQStatus
}

func (c *FAQCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.KbID > 0 {
		db.Where(tableName+".kb_id = ?", c.KbID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.Question != "" {
		db.Where(tableName+".question = ?", c.Question)
	}
	if c.Status != "" {
		db.Where(tableName+".status = ?", c.Status)
	}
}

type FAQDao struct {
	*gormdao.Dao[model.FAQEntity, model.FAQEntityList]
}

func NewFAQDao() *FAQDao {
	return &FAQDao{
		Dao: gormdao.NewDao[model.FAQEntity, model.FAQEntityList](
			model.TableNameFAQ, "FAQDao",
			dbclient.RagForgeDB,
		),
	}
}
