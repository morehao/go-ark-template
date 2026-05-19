package model

import (
	"gorm.io/gorm"
)

type SessionEntity struct {
	gorm.Model
	TenantID    uint   `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	UserID      uint   `gorm:"column:user_id;type:bigint unsigned;not null;default:0;index;comment:用户id"`
	KbID        uint   `gorm:"column:kb_id;type:bigint unsigned;not null;default:0;index;comment:知识库id"`
	Title       string `gorm:"column:title;type:varchar(500);not null;default:'';comment:会话标题"`
	Description string `gorm:"column:description;type:varchar(1000);not null;default:'';comment:会话描述"`
	IsPinned    bool   `gorm:"column:is_pinned;type:tinyint(1);not null;default:0;comment:是否置顶"`
}

type SessionEntityList []SessionEntity

func (SessionEntity) TableName() string {
	return TableNameSession
}

