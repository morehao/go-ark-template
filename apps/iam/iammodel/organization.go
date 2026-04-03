package iammodel

import (
	"gorm.io/gorm"
)

// OrganizationEntity 产品线管理表结构体
type OrganizationEntity struct {
	gorm.Model
	CreatedBy        uint   `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy        uint   `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Domain           string `gorm:"column:domain;type:varchar(255);;default '';comment: 产品线域名"`
	Logo             string `gorm:"column:logo;type:varchar(255);;default '';comment: 产品线Logo"`
	Description      string `gorm:"column:description;type:varchar(255);;default '';comment: 产品线描述"`
	SortOrder        int32  `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	Status           string `gorm:"column:status;type:varchar(16);;default active;comment: 状态: active-正常 inactive-停用"`
	OrganizationCode string `gorm:"column:organization_code;type:varchar(32);not null;default '';comment: 产品线编码"`
	OrganizationName string `gorm:"column:organization_name;type:varchar(64);not null;default '';comment: 产品线名称"`
	UpdatedBy        uint   `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
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
