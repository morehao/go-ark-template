package iammodel

import (
	"gorm.io/gorm"
)

type TenantStatus string

const (
	TenantStatusEnabled  TenantStatus = "enabled"
	TenantStatusTrial    TenantStatus = "trial"
	TenantStatusExpired  TenantStatus = "expired"
	TenantStatusDisabled TenantStatus = "disabled"
)

type TenantEntity struct {
	gorm.Model
	Address                 string       `gorm:"column:address;type:varchar(255);;default '';comment: 租户地址"`
	ContactEmail            string       `gorm:"column:contact_email;type:varchar(64);;default '';comment: 联系邮箱"`
	ContactPhone            string       `gorm:"column:contact_phone;type:varchar(16);;default '';comment: 联系电话"`
	CreatedBy               uint         `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy               uint         `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	LegalPerson             string       `gorm:"column:legal_person;type:varchar(32);;default '';comment: 法人代表"`
	Logo                    string       `gorm:"column:logo;type:varchar(255);;default '';comment: 租户Logo"`
	OrgID                   uint         `gorm:"column:org_id;type:bigint;not null;default 0;comment: 所属组织ID"`
	ShortName               string       `gorm:"column:short_name;type:varchar(64);;default '';comment: 租户简称"`
	Status                  TenantStatus `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 trial-试用 expired-已过期 disabled-停用"`
	TenantCode              string       `gorm:"column:tenant_code;type:varchar(32);not null;default '';comment: 租户编码"`
	TenantName              string       `gorm:"column:tenant_name;type:varchar(128);not null;default '';comment: 租户名称"`
	UnifiedSocialCreditCode string       `gorm:"column:unified_social_credit_code;type:varchar(18);;default '';comment: 统一社会信用代码(18位)"`
	UpdatedBy               uint         `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
}

type TenantEntityList []TenantEntity

const TableNameTenant = "iam_tenant"

func (TenantEntity) TableName() string {
	return TableNameTenant
}

func (l TenantEntityList) ToMap() map[uint]TenantEntity {
	m := make(map[uint]TenantEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
