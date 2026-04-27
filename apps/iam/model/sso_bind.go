package model

import (
	"gorm.io/gorm"
)

type SSOType string

const (
	SSOTypeWechat SSOType = "wechat"
	SSOTypeOIDC   SSOType = "oidc"
)

type SSOBindEntity struct {
	gorm.Model
	OrgID    uint   `gorm:"column:org_id;type:bigint;not null;default 0;comment: 组织ID"`
	TenantID uint   `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 租户ID"`
	UserID   uint   `gorm:"column:user_id;type:bigint;not null;default 0;comment: 用户ID"`
	SSOType  string `gorm:"column:sso_type;type:varchar(32);not null;default '';comment: SSO类型，如wechat/oidc"`
	OpenID   string `gorm:"column:open_id;type:varchar(128);not null;default '';comment: OpenID"`
}

type SSOBindEntityList []SSOBindEntity

const TableNameSSOBind = "iam_sso_bind"

func (SSOBindEntity) TableName() string {
	return TableNameSSOBind
}

func (l SSOBindEntityList) ToMap() map[uint]SSOBindEntity {
	m := make(map[uint]SSOBindEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}