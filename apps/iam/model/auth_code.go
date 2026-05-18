package model

import (
	"time"

	"gorm.io/gorm"
)

type AuthCodeEntity struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Code                string         `gorm:"column:code;type:varchar(64);unique;not null;comment:授权码" json:"code"`
	ClientID            string         `gorm:"column:client_id;type:varchar(64);not null;comment:Client ID" json:"client_id"`
	PersonID            uint           `gorm:"column:person_id;type:bigint;not null;default:0;comment:自然人ID" json:"person_id"`
	TenantID            uint           `gorm:"column:tenant_id;type:bigint;not null;default:0;comment:租户ID" json:"tenant_id"`
	OrgID               uint           `gorm:"column:org_id;type:bigint;not null;default:0;comment:组织ID" json:"org_id"`
	RedirectURI         string         `gorm:"column:redirect_uri;type:varchar(255);not null;comment:重定向URI" json:"redirect_uri"`
	Scope               string         `gorm:"column:scope;type:varchar(255);default:openid,profile;comment:请求的scope" json:"scope"`
	State               string         `gorm:"column:state;type:varchar(128);comment:state参数，防CSRF" json:"state"`
	CodeChallenge       string         `gorm:"column:code_challenge;type:varchar(64);comment:PKCE code_challenge" json:"code_challenge"`
	CodeChallengeMethod string         `gorm:"column:code_challenge_method;type:varchar(8);default:S256;comment:PKCE challenge方法" json:"code_challenge_method"`
	ExpiresAt           time.Time      `gorm:"column:expires_at;type:datetime(3);not null;comment:过期时间" json:"expires_at"`
	Used                bool           `gorm:"column:used;type:tinyint(1);not null;default:0;comment:是否已使用" json:"used"`
	UsedAt              *time.Time     `gorm:"column:used_at;type:datetime(3);comment:使用时间" json:"used_at,omitempty"`
	CreatedAt           time.Time      `gorm:"column:created_at;type:datetime(3);not null" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;type:datetime(3);not null" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);index" json:"deleted_at,omitempty"`
}

const TableNameAuthCode = "iam_auth_code"

func (AuthCodeEntity) TableName() string {
	return TableNameAuthCode
}
