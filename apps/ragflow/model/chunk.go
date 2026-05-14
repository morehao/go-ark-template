package model

import (
	"gorm.io/gorm"
)

type ChunkEntity struct {
	gorm.Model
	DocumentID uint   `gorm:"column:document_id;type:bigint unsigned;not null;default 0;comment: 文档id"`
	Content    string `gorm:"column:content;type:text;not null;comment: 内容"`
	Vector     []byte `gorm:"column:vector;type:blob;comment: 向量"`
	Metadata   string `gorm:"column:metadata;type:text;comment: 元数据"`
}

type ChunkEntityList []ChunkEntity

const TableNameChunk = "ragflow_chunk"

func (ChunkEntity) TableName() string {
	return TableNameChunk
}

func (l ChunkEntityList) ToMap() map[uint]ChunkEntity {
	m := make(map[uint]ChunkEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}