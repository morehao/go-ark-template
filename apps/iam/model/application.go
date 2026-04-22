package model

import (
	"gorm.io/gorm"
)

// ApplicationEntity 应用管理表结构体
type ApplicationEntity struct {
	gorm.Model
	AppCode     string `gorm:"column:app_code;type:varchar(32);not null;default '';comment: 应用编码"`
	AppName     string `gorm:"column:app_name;type:varchar(64);not null;default '';comment: 应用名称"`
	AppType     string `gorm:"column:app_type;type:varchar(16);;default web;comment: 应用类型: web-网页 app-移动端 mini-小程序"`
	CallbackUrl string `gorm:"column:callback_url;type:varchar(255);;default '';comment: 回调URL"`
	CreatedBy   uint   `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy   uint   `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Description string `gorm:"column:description;type:varchar(255);;default '';comment: 应用描述"`
	HomepageUrl string `gorm:"column:homepage_url;type:varchar(255);;default '';comment: 应用首页URL"`
	Logo        string `gorm:"column:logo;type:varchar(255);;default '';comment: 应用Logo"`
	SortOrder   int32  `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	Status      string `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	UpdatedBy   uint   `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
}

type ApplicationEntityList []ApplicationEntity

const TableNameApplication = "iam_application"

func (ApplicationEntity) TableName() string {
	return TableNameApplication
}

func (l ApplicationEntityList) ToMap() map[uint]ApplicationEntity {
	m := make(map[uint]ApplicationEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
