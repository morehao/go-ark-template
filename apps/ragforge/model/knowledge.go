package model

import (
	"gorm.io/gorm"
)

type KnowledgeType string

const (
	KnowledgeTypeFile    KnowledgeType = "file"
	KnowledgeTypeURL     KnowledgeType = "url"
	KnowledgeTypeManual  KnowledgeType = "manual"
)

type ParseStatus string

const (
	ParseStatusPending   ParseStatus = "pending"
	ParseStatusParsing   ParseStatus = "parsing"
	ParseStatusCompleted ParseStatus = "completed"
	ParseStatusFailed    ParseStatus = "failed"
)

type KnowledgeEntity struct {
	gorm.Model
	KbID        uint         `gorm:"column:kb_id;type:bigint unsigned;not null;default:0;index;comment:知识库id"`
	TenantID    uint         `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Type        KnowledgeType `gorm:"column:type;type:varchar(50);not null;default:'file';comment:知识类型"`
	Title       string       `gorm:"column:title;type:varchar(500);not null;default:'';comment:标题"`
	Content     string       `gorm:"column:content;type:longtext;not null;comment:内容"`
	FileURL     string       `gorm:"column:file_url;type:varchar(500);not null;default:'';comment:文件地址"`
	SourceURL   string       `gorm:"column:source_url;type:varchar(500);not null;default:'';comment:来源地址"`
	ParseStatus ParseStatus  `gorm:"column:parse_status;type:varchar(50);not null;default:'pending';comment:解析状态"`
	FileSize    int64        `gorm:"column:file_size;type:bigint;not null;default:0;comment:文件大小(bytes)"`
	CreatorID   uint         `gorm:"column:creator_id;type:bigint unsigned;not null;default:0;index;comment:创建人id"`
}

type KnowledgeEntityList []KnowledgeEntity

func (KnowledgeEntity) TableName() string {
	return TableNameKnowledge
}

