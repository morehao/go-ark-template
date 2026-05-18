package model

import (
	"time"

	"gorm.io/gorm"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
	TokenTypeID      TokenType = "id"
)

type TokenEntity struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenID          string         `gorm:"column:token_id;type:varchar(64);unique;not null;comment:Token唯一标识" json:"token_id"`
	PersonID         uint           `gorm:"column:person_id;type:bigint;not null;default:0;comment:自然人ID" json:"person_id"`
	UserID           uint           `gorm:"column:user_id;type:bigint;not null;default:0;comment:用户ID" json:"user_id"`
	ClientID         string         `gorm:"column:client_id;type:varchar(64);not null;comment:Client ID" json:"client_id"`
	TenantID         uint           `gorm:"column:tenant_id;type:bigint;not null;default:0;comment:租户ID" json:"tenant_id"`
	OrgID            uint           `gorm:"column:org_id;type:bigint;not null;default:0;comment:组织ID" json:"org_id"`
	TokenType        TokenType      `gorm:"column:token_type;type:varchar(16);not null;comment:Token类型" json:"token_type"`
	AccessTokenHash  string         `gorm:"column:access_token_hash;type:varchar(128);comment:Access Token哈希" json:"access_token_hash,omitempty"`
	RefreshTokenHash string         `gorm:"column:refresh_token_hash;type:varchar(128);comment:Refresh Token哈希" json:"refresh_token_hash,omitempty"`
	Scopes           string         `gorm:"column:scopes;type:varchar(255);default:openid,profile;comment:授权的scopes" json:"scopes"`
	ExpiresAt        time.Time      `gorm:"column:expires_at;type:datetime(3);not null;comment:过期时间" json:"expires_at"`
	Revoked          bool           `gorm:"column:revoked;type:tinyint(1);not null;default:0;comment:是否撤销" json:"revoked"`
	RevokedAt        *time.Time     `gorm:"column:revoked_at;type:datetime(3);comment:撤销时间" json:"revoked_at,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at;type:datetime(3);not null" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:datetime(3);not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);index" json:"deleted_at,omitempty"`
}

const TableNameToken = "iam_token"

func (TokenEntity) TableName() string {
	return TableNameToken
}