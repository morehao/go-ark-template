package model

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type ChunkEntity struct {
	gorm.Model
	KnowledgeID uint            `gorm:"column:knowledge_id;type:bigint unsigned;not null;default:0;index;comment:知识id"`
	KbID        uint            `gorm:"column:kb_id;type:bigint unsigned;not null;default:0;index;comment:知识库id"`
	TenantID    uint            `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Content     string          `gorm:"column:content;type:longtext;not null;comment:块内容"`
	SeqID       int             `gorm:"column:seq_id;type:int;not null;default:0;comment:序号"`
	Tokens      int             `gorm:"column:tokens;type:int;not null;default:0;comment:token数量"`
	Vector      pgvector.Vector `gorm:"type:vector(1536);column:vector;comment:向量"`
	MetaInfo    string          `gorm:"column:meta_info;type:jsonb;not null;default:'{}';comment:元信息"`
}

type ChunkEntityList []ChunkEntity

func (ChunkEntity) TableName() string {
	return TableNameChunk
}

