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
	Sequence   int32  `gorm:"column:sequence;type:int;;default 0;comment: 排序"`
	Status      string `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	UpdatedBy   uint   `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	ClientID       string `gorm:"column:client_id;type:varchar(64);unique;comment:OAuth2 Client ID"`
	ClientSecret   string `gorm:"column:client_secret;type:varchar(255);comment:OAuth2 Client Secret"`
	ClientType     string `gorm:"column:client_type;type:varchar(16);default:web;comment:客户端类型: web/app/spa/mini"`
	PkceRequired   bool   `gorm:"column:pkce_required;type:tinyint(1);default:0;comment:是否强制PKCE"`
	AllowedScopes  string `gorm:"column:allowed_scopes;type:varchar(255);default:openid,profile;comment:允许的scopes"`
	AllowedCallbacks string `gorm:"column:allowed_callbacks;type:text;comment:允许的重定向URI，JSON数组"`
}

type ApplicationEntityList []ApplicationEntity

const TableNameApplication = "iam_application"

type AppStatus string

const (
	AppStatusEnabled  AppStatus = "enabled"
	AppStatusDisabled AppStatus = "disabled"
)

type ClientType string

const (
	ClientTypeWeb  ClientType = "web"
	ClientTypeApp  ClientType = "app"
	ClientTypeSpa  ClientType = "spa"
	ClientTypeMini ClientType = "mini"
)

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
