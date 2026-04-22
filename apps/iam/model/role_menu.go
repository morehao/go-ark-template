package model

import (
	"gorm.io/gorm"
)

// RoleMenuEntity 角色菜单关联表结构体
type RoleMenuEntity struct {
	gorm.Model
	TenantID  uint `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 租户ID(冗余)"`
	CreatedBy uint `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy uint `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	MenuID    uint `gorm:"column:menu_id;type:bigint;not null;default '';comment: 菜单ID"`
	RoleID    uint `gorm:"column:role_id;type:bigint;not null;default '';comment: 角色ID"`
	UpdatedBy uint `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
}

type RoleMenuEntityList []RoleMenuEntity

const TableNameRoleMenu = "iam_role_menu"

func (RoleMenuEntity) TableName() string {
	return TableNameRoleMenu
}

func (l RoleMenuEntityList) ToMap() map[uint]RoleMenuEntity {
	m := make(map[uint]RoleMenuEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
