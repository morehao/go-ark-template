package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type TagCond struct {
	*gormdao.BaseCond
	KbID     uint
	TenantID uint
}

func (c *TagCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.KbID > 0 {
		db.Where(tableName+".kb_id = ?", c.KbID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
}

type TagDao struct {
	*gormdao.Dao[model.TagEntity, model.TagEntityList]
}

func NewTagDao() *TagDao {
	return &TagDao{
		Dao: gormdao.NewDao[model.TagEntity, model.TagEntityList](
			model.TableNameTag, "TagDao",
			dbclient.RagForgeDB,
		),
	}
}
