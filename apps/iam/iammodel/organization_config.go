package iammodel

import (
	"gorm.io/gorm"
)

// OrganizationConfigEntity 产品线配置表结构体
type OrganizationConfigEntity struct {
	gorm.Model
	ConfigGroup    string `gorm:"column:config_group;type:varchar(32);;default general;comment: 配置分组: general-通用/auth-认证/theme-主题等"`
	ConfigKey      string `gorm:"column:config_key;type:varchar(100);not null;default '';comment: 配置键"`
	ConfigType     string `gorm:"column:config_type;type:varchar(32);;default string;comment: 配置类型: string/json/boolean/number"`
	ConfigValue    string `gorm:"column:config_value;type:text;;default '';comment: 配置值(支持JSON)"`
	Description    string `gorm:"column:description;type:varchar(255);not null;default '';comment: 配置说明"`
	SortOrder      int32  `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	OrganizationID uint   `gorm:"column:organization_id;type:bigint;not null;default 0;comment: 产品线ID"`
}

type OrganizationConfigEntityList []OrganizationConfigEntity

const TableNameOrganizationConfig = "iam_organization_config"

func (OrganizationConfigEntity) TableName() string {
	return TableNameOrganizationConfig
}

func (l OrganizationConfigEntityList) ToMap() map[uint]OrganizationConfigEntity {
	m := make(map[uint]OrganizationConfigEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
