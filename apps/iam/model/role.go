package model

import (
	"gorm.io/gorm"
)

type RoleDataScope string

const (
	RoleDataScopeAll        RoleDataScope = "all"
	RoleDataScopeDeptAndSub RoleDataScope = "dept_and_sub"
	RoleDataScopeDept       RoleDataScope = "dept"
	RoleDataScopeSelf       RoleDataScope = "self"
	RoleDataScopeCustom     RoleDataScope = "custom"
)

type RoleType string

const (
	RoleTypeCustom RoleType = "custom"
	RoleTypeSystem RoleType = "system"
)

type RoleStatus string

const (
	RoleStatusEnabled  RoleStatus = "enabled"
	RoleStatusDisabled RoleStatus = "disabled"
)

// RoleEntity 角色管理表结构体
type RoleEntity struct {
	gorm.Model
	TenantID    uint          `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 所属租户ID"`
	CreatedBy   uint          `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DataScope   RoleDataScope `gorm:"column:data_scope;type:varchar(16);;default all;comment: 数据权限范围: all-全部 dept_and_sub-本部门及以下 dept-本部门 self-仅本人 custom-自定义"`
	DeletedBy   uint          `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Description string        `gorm:"column:description;type:varchar(255);;default '';comment: 角色描述"`
	RoleCode    string        `gorm:"column:role_code;type:varchar(32);not null;default '';comment: 角色编码"`
	RoleName    string        `gorm:"column:role_name;type:varchar(64);not null;default '';comment: 角色名称"`
	RoleType    RoleType      `gorm:"column:role_type;type:varchar(16);;default custom;comment: 角色类型: custom-自定义 system-系统内置"`
	Sequence   int32         `gorm:"column:sequence;type:int;;default 0;comment: 排序"`
	Status      RoleStatus    `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	UpdatedBy   uint          `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
}

type RoleEntityList []RoleEntity

const TableNameRole = "iam_role"

func (RoleEntity) TableName() string {
	return TableNameRole
}

func (l RoleEntityList) ToMap() map[uint]RoleEntity {
	m := make(map[uint]RoleEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
