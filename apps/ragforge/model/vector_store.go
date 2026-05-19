package model

import (
	"gorm.io/gorm"
)

type EngineType string

const (
	EngineTypeElasticsearch EngineType = "elasticsearch"
	EngineTypeMilvus       EngineType = "milvus"
	EngineTypePgvector     EngineType = "pgvector"
)

type VectorStoreStatus string

const (
	VectorStoreStatusActive   VectorStoreStatus = "active"
	VectorStoreStatusInactive VectorStoreStatus = "inactive"
)

type VectorStoreEntity struct {
	gorm.Model
	TenantID   uint              `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Name       string            `gorm:"column:name;type:varchar(255);not null;default:'';comment:向量库名称"`
	EngineType EngineType        `gorm:"column:engine_type;type:varchar(50);not null;default:'elasticsearch';comment:引擎类型"`
	Config     string            `gorm:"column:config;type:jsonb;not null;default:'{}';comment:连接配置"`
	Status     VectorStoreStatus `gorm:"column:status;type:varchar(50);not null;default:'active';comment:状态"`
}

type VectorStoreEntityList []VectorStoreEntity

func (VectorStoreEntity) TableName() string {
	return TableNameVectorStore
}

