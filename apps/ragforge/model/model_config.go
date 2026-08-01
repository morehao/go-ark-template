package model

import (
	"gorm.io/gorm"
)

type ModelType string

const (
	ModelTypeChat      ModelType = "chat"
	ModelTypeEmbedding ModelType = "embedding"
	ModelTypeRerank    ModelType = "rerank"
	ModelTypeVLM       ModelType = "vlm"
	ModelTypeASR       ModelType = "asr"
)

type ModelStatus string

const (
	ModelStatusActive   ModelStatus = "active"
	ModelStatusInactive ModelStatus = "inactive"
)

type ModelEntity struct {
	gorm.Model
	TenantID  uint        `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Name      string      `gorm:"column:name;type:varchar(255);not null;default:'';comment:模型名称"`
	ModelType ModelType   `gorm:"column:model_type;type:varchar(50);not null;default:'chat';comment:模型类型"`
	Provider  string      `gorm:"column:provider;type:varchar(100);not null;default:'';comment:提供商"`
	ModelName string      `gorm:"column:model_name;type:varchar(255);not null;default:'';comment:模型标识"`
	Config    string      `gorm:"column:config;type:jsonb;not null;default:'{}';comment:模型配置"`
	Status    ModelStatus `gorm:"column:status;type:varchar(50);not null;default:'active';comment:状态"`
}

type ModelEntityList []ModelEntity

func (ModelEntity) TableName() string {
	return TableNameModel
}

