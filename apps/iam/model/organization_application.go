package model

import (
	"gorm.io/gorm"
)

// OrganizationApplicationEntity 组织应用关系表结构体
type OrganizationApplicationEntity struct {
	gorm.Model
	AppID     uint `gorm:"column:app_id;type:bigint;not null;default '';comment: 应用ID"`
	CreatedBy uint `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	OrgID     uint `gorm:"column:org_id;type:bigint;not null;default '';comment: 组织ID"`
}

type OrganizationApplicationEntityList []OrganizationApplicationEntity

const TableNameOrganizationApplication = "iam_organization_application"

func (OrganizationApplicationEntity) TableName() string {
	return TableNameOrganizationApplication
}

func (l OrganizationApplicationEntityList) ToMap() map[uint]OrganizationApplicationEntity {
	m := make(map[uint]OrganizationApplicationEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
