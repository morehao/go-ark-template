package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/dbaccess/gormdao"
	"gorm.io/gorm"
)

type SessionCond struct {
	*gormdao.BaseCond
	TenantID uint
	UserID   uint
	KbID     uint
}

func (c *SessionCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.UserID > 0 {
		db.Where(tableName+".user_id = ?", c.UserID)
	}
	if c.KbID > 0 {
		db.Where(tableName+".kb_id = ?", c.KbID)
	}
}

type SessionDao struct {
	*gormdao.Dao[model.SessionEntity, model.SessionEntityList]
}

func NewSessionDao() *SessionDao {
	return &SessionDao{
		Dao: gormdao.NewDao[model.SessionEntity, model.SessionEntityList](
			model.TableNameSession, "SessionDao",
			dbclient.RagForgeDB,
		),
	}
}
