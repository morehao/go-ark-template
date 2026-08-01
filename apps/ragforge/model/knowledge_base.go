package model

import (
	"gorm.io/gorm"
)

type KBType string

const (
	KBTypeNormal  KBType = "normal"
	KBTypeFAQ     KBType = "faq"
)

type ParserEngine string

const (
	ParserEngineDefault  ParserEngine = "default"
	ParserEngineDocling  ParserEngine = "docling"
	ParserEngineTika     ParserEngine = "tika"
)

type KnowledgeBaseEntity struct {
	gorm.Model
	TenantID        uint         `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Name            string       `gorm:"column:name;type:varchar(255);not null;default:'';comment:知识库名称"`
	Description     string       `gorm:"column:description;type:varchar(500);not null;default:'';comment:知识库描述"`
	KBType          KBType       `gorm:"column:kb_type;type:varchar(50);not null;default:'normal';comment:知识库类型"`
	ParserEngine    ParserEngine `gorm:"column:parser_engine;type:varchar(50);not null;default:'default';comment:解析引擎"`
	EmbeddingConfig string       `gorm:"column:embedding_config;type:jsonb;not null;default:'{}';comment:嵌入配置"`
	IndexStrategy   string       `gorm:"column:index_strategy;type:jsonb;not null;default:'{}';comment:索引策略"`
	CreatorID       uint         `gorm:"column:creator_id;type:bigint unsigned;not null;default:0;index;comment:创建人id"`
}

type KnowledgeBaseEntityList []KnowledgeBaseEntity

func (KnowledgeBaseEntity) TableName() string {
	return TableNameKnowledgeBase
}

