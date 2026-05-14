package model

import (
	"gorm.io/gorm"
)

type DocumentEntity struct {
	gorm.Model
	KnowledgeBaseID uint   `gorm:"column:knowledge_base_id;type:bigint unsigned;not null;default 0;comment: 知识库id"`
	Name            string `gorm:"column:name;type:varchar(255);not null;default '';comment: 文档名称"`
	Type            string `gorm:"column:type;type:varchar(32);not null;default '';comment: 文档类型"`
	Location        string `gorm:"column:location;type:text;comment: 文档位置"`
	Size            int64  `gorm:"column:size;type:bigint;not null;default 0;comment: 文档大小"`
	Status          string `gorm:"column:status;type:varchar(32);not null;default 'pending';comment: 状态"`
	ChunkStatus     string `gorm:"column:chunk_status;type:varchar(32);not null;default 'pending';comment: 分块状态"`
	VectorStatus    string `gorm:"column:vector_status;type:varchar(32);not null;default 'pending';comment: 向量化状态"`
}

type DocumentEntityList []DocumentEntity

const TableNameDocument = "ragflow_document"

func (DocumentEntity) TableName() string {
	return TableNameDocument
}

func (l DocumentEntityList) ToMap() map[uint]DocumentEntity {
	m := make(map[uint]DocumentEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}