package model

import (
	"gorm.io/gorm"
)

type OrgStatus string

const (
	OrgStatusEnabled  OrgStatus = "enabled"
	OrgStatusDisabled OrgStatus = "disabled"
)

// OrganizationEntity 组织管理表结构体
type OrganizationEntity struct {
	gorm.Model
	CreatedBy   uint      `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy   uint      `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Domain      string    `gorm:"column:domain;type:varchar(255);;default '';comment: 组织域名"`
	Logo        string    `gorm:"column:logo;type:varchar(255);;default '';comment: 组织Logo"`
	Description string    `gorm:"column:description;type:varchar(255);;default '';comment: 组织描述"`
	SortOrder   int32     `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	Status      OrgStatus `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	OrgCode     string    `gorm:"column:org_code;type:varchar(32);not null;default '';comment: 组织编码"`
	OrgName     string    `gorm:"column:org_name;type:varchar(64);not null;default '';comment: 组织名称"`
	UpdatedBy   uint      `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
}

type OrganizationEntityList []OrganizationEntity

const TableNameOrganization = "iam_organization"

func (OrganizationEntity) TableName() string {
	return TableNameOrganization
}

func (l OrganizationEntityList) ToMap() map[uint]OrganizationEntity {
	m := make(map[uint]OrganizationEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
