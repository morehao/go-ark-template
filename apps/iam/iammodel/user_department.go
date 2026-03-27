package iammodel

import (
	"gorm.io/gorm"
)

// UserDepartmentEntity 用户部门关系表结构体
type UserDepartmentEntity struct {
	gorm.Model
	TenantID  uint   `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 租户ID(冗余)"`
	CreatedBy uint   `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy uint   `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	DeptID    uint   `gorm:"column:dept_id;type:bigint;not null;default '';comment: 部门ID"`
	DeptType  string `gorm:"column:dept_type;type:varchar(16);;default primary;comment: 部门类型: primary-主部门 secondary-其他部门"`
	UpdatedBy uint   `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	UserID    uint   `gorm:"column:user_id;type:bigint;not null;default '';comment: 用户ID"`
}

type UserDepartmentEntityList []UserDepartmentEntity

const TableNameUserDepartment = "iam_user_department"

func (UserDepartmentEntity) TableName() string {
	return TableNameUserDepartment
}

func (l UserDepartmentEntityList) ToMap() map[uint]UserDepartmentEntity {
	m := make(map[uint]UserDepartmentEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
