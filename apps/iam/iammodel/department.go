package iammodel

import (
	"gorm.io/gorm"
)

type DeptStatus string

const (
	DeptStatusEnabled  DeptStatus = "enabled"
	DeptStatusDisabled DeptStatus = "disabled"
)

// DepartmentEntity 部门管理表结构体
type DepartmentEntity struct {
	gorm.Model
	TenantID  uint       `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 所属租户ID"`
	CreatedBy uint       `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy uint       `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	DeptCode  string     `gorm:"column:dept_code;type:varchar(32);not null;default '';comment: 部门编码"`
	DeptLevel int32      `gorm:"column:dept_level;type:int;;default 1;comment: 部门层级"`
	DeptName  string     `gorm:"column:dept_name;type:varchar(64);not null;default '';comment: 部门名称"`
	DeptPath  string     `gorm:"column:dept_path;type:varchar(512);;default '';comment: 部门路径: /1/2/3/"`
	LeaderID  uint       `gorm:"column:leader_id;type:bigint;;default '';comment: 部门负责人ID"`
	ParentID  uint       `gorm:"column:parent_id;type:bigint;;default 0;comment: 父部门ID,0表示根部门"`
	SortOrder int32      `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	Status    DeptStatus `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	UpdatedBy uint       `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
}

type DepartmentEntityList []DepartmentEntity

const TableNameDepartment = "iam_department"

func (DepartmentEntity) TableName() string {
	return TableNameDepartment
}

func (l DepartmentEntityList) ToMap() map[uint]DepartmentEntity {
	m := make(map[uint]DepartmentEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
