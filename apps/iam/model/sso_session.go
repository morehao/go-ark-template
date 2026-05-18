package model

import (
	"time"

	"gorm.io/gorm"
)

type SsoSessionEntity struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID      string         `gorm:"column:session_id;type:varchar(64);unique;not null;comment:SSO会话ID" json:"session_id"`
	PersonID       uint           `gorm:"column:person_id;type:bigint;not null;default:0;comment:自然人ID" json:"person_id"`
	OrgID          uint           `gorm:"column:org_id;type:bigint;not null;default:0;comment:组织ID" json:"org_id"`
	LoginTime      time.Time      `gorm:"column:login_time;type:datetime(3);not null;comment:登录时间" json:"login_time"`
	LastActiveTime time.Time      `gorm:"column:last_active_time;type:datetime(3);not null;comment:最后活跃时间" json:"last_active_time"`
	ExpiresAt      time.Time      `gorm:"column:expires_at;type:datetime(3);not null;comment:过期时间" json:"expires_at"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:datetime(3);not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:datetime(3);not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);index" json:"deleted_at,omitempty"`
}

const TableNameSsoSession = "iam_sso_session"

func (SsoSessionEntity) TableName() string {
	return TableNameSsoSession
}