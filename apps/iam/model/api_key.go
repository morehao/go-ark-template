package model

import (
	"gorm.io/gorm"
)

type ApiKeyStatus string

const (
	ApiKeyStatusEnabled  ApiKeyStatus = "enabled"
	ApiKeyStatusDisabled ApiKeyStatus = "disabled"
)

type ApiKeyAccessPolicy string

const (
	ApiKeyAccessPolicyAll ApiKeyAccessPolicy = "all"
	ApiKeyAccessPolicyIP  ApiKeyAccessPolicy = "ip"
)

type ApiKeyEntity struct {
	ID                  uint           `gorm:"column:id;type:bigint;autoIncrement;primaryKey"`
	TenantID            uint           `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 租户ID"`
	UserID              uint           `gorm:"column:user_id;type:bigint;not null;default 0;comment: 关联用户ID"`
	AppID               uint           `gorm:"column:app_id;type:bigint;not null;default 0;comment: 应用ID"`
	KeyName             string         `gorm:"column:key_name;type:varchar(64);not null;comment: 密钥名称"`
	KeyPrefix           string         `gorm:"column:key_prefix;type:varchar(16);not null;comment: 密钥前缀"`
	PublicKey           string         `gorm:"column:public_key;type:text;not null;comment: 公钥"`
	EncryptedPrivateKey string         `gorm:"column:encrypted_private_key;type:text;not null;comment: 加密私钥"`
	AccessPolicy        ApiKeyAccessPolicy `gorm:"column:access_policy;type:varchar(16);default all;comment: 访问策略"`
	AllowedIPs          string         `gorm:"column:allowed_ips;type:text;;comment: 允许的IP列表(JSON)"`
	Scopes              string         `gorm:"column:scopes;type:varchar(255);;comment: 权限范围"`
	Status              ApiKeyStatus   `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	LastUsedAt          string         `gorm:"column:last_used_at;type:datetime(3);;comment: 最后使用时间"`
	ExpiresAt           string         `gorm:"column:expires_at;type:datetime(3);;comment: 过期时间"`
	CreatedAt           string         `gorm:"column:created_at;type:datetime(3);;comment: 创建时间"`
	UpdatedAt           string         `gorm:"column:updated_at;type:datetime(3);;comment: 更新时间"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);;comment: 删除时间"`
	CreatedBy           uint           `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	UpdatedBy           uint           `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	DeletedBy           uint           `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
}

type ApiKeyEntityList []ApiKeyEntity

const TableNameApiKey = "iam_api_key"

func (ApiKeyEntity) TableName() string {
	return TableNameApiKey
}

func (l ApiKeyEntityList) ToMap() map[uint]ApiKeyEntity {
	m := make(map[uint]ApiKeyEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
