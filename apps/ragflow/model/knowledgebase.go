package model

import (
	"gorm.io/gorm"
)

type KnowledgeBaseEntity struct {
	gorm.Model
	Name           string `gorm:"column:name;type:varchar(255);not null;default '';comment: 知识库名称"`
	Description    string `gorm:"column:description;type:text;comment: 知识库描述"`
	EmbeddingModel string `gorm:"column:embedding_model;type:varchar(128);not null;default '';comment:  embedding模型"`
	VectorStoreType string `gorm:"column:vector_store_type;type:varchar(32);not null;default '';comment: 向量库类型"`
	PermissionType  string `gorm:"column:permission_type;type:varchar(32);not null;default '';comment: 权限类型"`
	Status         string `gorm:"column:status;type:varchar(32);not null;default 'active';comment: 状态"`
	ChunkMethod     string `gorm:"column:chunk_method;type:varchar(32);not null;default '';comment: 分块方式"`
}

type KnowledgeBaseEntityList []KnowledgeBaseEntity

const TableNameKnowledgeBase = "ragflow_knowledgebase"

func (KnowledgeBaseEntity) TableName() string {
	return TableNameKnowledgeBase
}

func (l KnowledgeBaseEntityList) ToMap() map[uint]KnowledgeBaseEntity {
	m := make(map[uint]KnowledgeBaseEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}