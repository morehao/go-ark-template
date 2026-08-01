package model

import (
	"gorm.io/gorm"
)

type TagEntity struct {
	gorm.Model
	KbID     uint   `gorm:"column:kb_id;type:bigint unsigned;not null;default:0;index;comment:知识库id"`
	TenantID uint   `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Name     string `gorm:"column:name;type:varchar(100);not null;default:'';comment:标签名称"`
	Color    string `gorm:"column:color;type:varchar(50);not null;default:'';comment:标签颜色"`
}

type TagEntityList []TagEntity

func (TagEntity) TableName() string {
	return TableNameTag
}

