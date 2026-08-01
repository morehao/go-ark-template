package model

import (
	"gorm.io/gorm"
)

type TenantEntity struct {
	gorm.Model
	Name        string `gorm:"column:name;type:varchar(255);not null;default:'';comment:租户名称"`
	Code        string `gorm:"column:code;type:varchar(100);not null;uniqueIndex;default:'';comment:租户编码"`
	Description string `gorm:"column:description;type:varchar(500);not null;default:'';comment:租户描述"`
	Status      int32  `gorm:"column:status;type:int;not null;default:0;comment:状态 0-启用 1-停用"`
}

type TenantEntityList []TenantEntity

func (TenantEntity) TableName() string {
	return TableNameTenant
}

