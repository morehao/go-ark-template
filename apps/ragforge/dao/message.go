package dao

import (
	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/genericdao"
	"gorm.io/gorm"
)

type MessageCond struct {
	*genericdao.BaseCond
	SessionID uint
	TenantID  uint
	Role      model.MessageRole
	Keyword   string
}

func (c *MessageCond) BuildCondition(db *gorm.DB, tableName string) {
	if c.BaseCond != nil {
		c.BaseCond.BuildCondition(db, tableName)
	}
	if c.SessionID > 0 {
		db.Where(tableName+".session_id = ?", c.SessionID)
	}
	if c.TenantID > 0 {
		db.Where(tableName+".tenant_id = ?", c.TenantID)
	}
	if c.Role != "" {
		db.Where(tableName+".role = ?", c.Role)
	}
	if c.Keyword != "" {
		db.Where(tableName+".content LIKE ?", "%"+c.Keyword+"%")
	}
}

type MessageDao struct {
	*genericdao.GenericDao[model.MessageEntity, model.MessageEntityList]
}

func NewMessageDao() *MessageDao {
	return &MessageDao{
		GenericDao: genericdao.NewGenericDao[model.MessageEntity, model.MessageEntityList](
			model.TableNameMessage, "MessageDao",
			dbclient.RagForgeDB,
		),
	}
}
