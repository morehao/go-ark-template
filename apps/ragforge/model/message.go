package model

import (
	"gorm.io/gorm"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

type MessageEntity struct {
	gorm.Model
	SessionID  uint        `gorm:"column:session_id;type:bigint unsigned;not null;default:0;index;comment:会话id"`
	TenantID   uint        `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Role       MessageRole `gorm:"column:role;type:varchar(50);not null;default:'user';comment:角色"`
	Content    string      `gorm:"column:content;type:longtext;not null;comment:消息内容"`
	Metadata   string      `gorm:"column:metadata;type:jsonb;not null;default:'{}';comment:元数据"`
	TokenCount int         `gorm:"column:token_count;type:int;not null;default:0;comment:token数量"`
}

type MessageEntityList []MessageEntity

func (MessageEntity) TableName() string {
	return TableNameMessage
}

