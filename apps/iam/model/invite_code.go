package model

import (
	"time"

	"gorm.io/gorm"
)

type InviteCodeStatus string

const (
	InviteCodeStatusActive   InviteCodeStatus = "active"
	InviteCodeStatusDisabled InviteCodeStatus = "disabled"
	InviteCodeStatusExpired  InviteCodeStatus = "expired"
)

type InviteCodeEntity struct {
	gorm.Model
	OrgID       uint             `gorm:"column:org_id;type:bigint;not null;default 0;comment: 组织ID"`
	TenantID    uint             `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 租户ID"`
	Code        string           `gorm:"column:code;type:varchar(32);not null;default '';comment: 邀请码"`
	Status      InviteCodeStatus `gorm:"column:status;type:varchar(16);default 'active';comment: 状态"`
	ExpiredAt   *time.Time       `gorm:"column:expired_at;type:datetime;"`
	MaxUseCount int              `gorm:"column:max_use_count;type:int;default 0;comment: 最大使用次数，0表示不限制"`
	UseCount    int              `gorm:"column:use_count;type:int;default 0;comment: 已使用次数"`
	CreatedBy   uint             `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
}

type InviteCodeEntityList []InviteCodeEntity

const TableNameInviteCode = "iam_invite_code"

func (InviteCodeEntity) TableName() string {
	return TableNameInviteCode
}

func (l InviteCodeEntityList) ToMap() map[uint]InviteCodeEntity {
	m := make(map[uint]InviteCodeEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
