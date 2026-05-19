package model

import (
	"gorm.io/gorm"
)

type FAQStatus string

const (
	FAQStatusActive   FAQStatus = "active"
	FAQStatusInactive FAQStatus = "inactive"
)

type FAQEntity struct {
	gorm.Model
	KbID             uint      `gorm:"column:kb_id;type:bigint unsigned;not null;default:0;index;comment:知识库id"`
	TenantID         uint      `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Question         string    `gorm:"column:question;type:text;not null;comment:问题"`
	Answer           string    `gorm:"column:answer;type:longtext;not null;comment:答案"`
	SimilarQuestions string    `gorm:"column:similar_questions;type:jsonb;not null;default:'[]';comment:相似问题"`
	Tags             string    `gorm:"column:tags;type:jsonb;not null;default:'[]';comment:标签"`
	Status           FAQStatus `gorm:"column:status;type:varchar(50);not null;default:'active';comment:状态"`
	CreatorID        uint      `gorm:"column:creator_id;type:bigint unsigned;not null;default:0;index;comment:创建人id"`
}

type FAQEntityList []FAQEntity

func (FAQEntity) TableName() string {
	return TableNameFAQ
}

