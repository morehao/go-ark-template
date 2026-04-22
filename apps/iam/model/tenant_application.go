package model

import (
	"gorm.io/gorm"
)

// TenantApplicationEntity 租户应用关系表结构体
type TenantApplicationEntity struct {
	gorm.Model
	AppID     uint `gorm:"column:app_id;type:bigint;not null;default '';comment: 应用ID"`
	CreatedBy uint `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	TenantID  uint `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 租户ID"`
}

type TenantApplicationEntityList []TenantApplicationEntity

const TableNameTenantApplication = "iam_tenant_application"

func (TenantApplicationEntity) TableName() string {
	return TableNameTenantApplication
}

func (l TenantApplicationEntityList) ToMap() map[uint]TenantApplicationEntity {
	m := make(map[uint]TenantApplicationEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
