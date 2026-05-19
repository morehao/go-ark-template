package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type FAQCond struct {
	*genericdao.BaseCond
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
	*genericdao.GenericDao[model.FAQEntity, model.FAQEntityList]
}

func NewFAQDao() *FAQDao {
	return &FAQDao{
		GenericDao: genericdao.NewGenericDao[model.FAQEntity, model.FAQEntityList](
			model.TableNameFAQ, "FAQDao",
			dbclient.RagForgeDB,
		),
	}
}
