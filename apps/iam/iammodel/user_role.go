package iammodel

import (
	"gorm.io/gorm"
)

// UserRoleEntity 用户角色关联表结构体
type UserRoleEntity struct {
	gorm.Model
	TenantID  uint `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 租户ID(冗余)"`
	CreatedBy uint `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy uint `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	RoleID    uint `gorm:"column:role_id;type:bigint;not null;default '';comment: 角色ID"`
	UpdatedBy uint `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	UserID    uint `gorm:"column:user_id;type:bigint;not null;default '';comment: 用户ID"`
}

type UserRoleEntityList []UserRoleEntity

const TableNameUserRole = "iam_user_role"

func (UserRoleEntity) TableName() string {
	return TableNameUserRole
}

func (l UserRoleEntityList) ToMap() map[uint]UserRoleEntity {
	m := make(map[uint]UserRoleEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
