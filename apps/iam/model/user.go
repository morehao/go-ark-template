package model

import (
	"time"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusEnabled  UserStatus = "enabled"
	UserStatusLocked   UserStatus = "locked"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusPending  UserStatus = "pending"
)

type UserType string

const (
	UserTypeNormal        UserType = "normal"
	UserTypeTenantAdmin   UserType = "tenant_admin"
	UserTypePlatformAdmin UserType = "platform_admin"
)

// UserEntity 用户管理表结构体
type UserEntity struct {
	gorm.Model
	TenantID    uint       `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 所属租户ID"`
	CreatedBy   uint       `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy   uint       `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	DeptID      uint       `gorm:"column:dept_id;type:bigint;;default '';comment: 主部门ID(冗余字段,方便查询,实际关联关系在iam_user_department表)"`
	EmployeeNo  string     `gorm:"column:employee_no;type:varchar(32);;default '';comment: 工号"`
	EntryDate   *time.Time `gorm:"column:entry_date;type:date;;default '';comment: 入职日期"`
	JobLevel    string     `gorm:"column:job_level;type:varchar(32);;default '';comment: 职级"`
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:datetime(3);;default NULL;comment: 最后登录时间"`
	LastLoginIp string     `gorm:"column:last_login_ip;type:varchar(45);;default '';comment: 最后登录IP(支持IPv6)"`
	LoginCount     int        `gorm:"column:login_count;type:int;;default 0;comment: 登录次数"`
	LoginFailCount int        `gorm:"column:login_fail_count;type:int;;default 0;comment: 连续登录失败次数"`
	LockedUntil    *time.Time `gorm:"column:locked_until;type:datetime(3);;default NULL;comment: 账户锁定截止时间"`
	PersonID       uint       `gorm:"column:person_id;type:bigint;not null;default '';comment: 自然人ID"`
	Position    string     `gorm:"column:position;type:varchar(64);;default '';comment: 职位"`
	Status      UserStatus `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-正常 locked-锁定 disabled-禁用"`
	UpdatedBy   uint       `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	UserType    UserType   `gorm:"column:user_type;type:varchar(16);;default normal;comment: 用户类型: normal-普通用户 tenant_admin-租户管理员 platform_admin-平台管理员(可管理所有租户)"`
	Username    string     `gorm:"column:username;type:varchar(32);not null;default '';comment: 用户名(租户用户:租户内唯一,平台管理员:全局唯一,需应用层保证)"`
}

type UserEntityList []UserEntity

const TableNameUser = "iam_user"

func (UserEntity) TableName() string {
	return TableNameUser
}

func (l UserEntityList) ToMap() map[uint]UserEntity {
	m := make(map[uint]UserEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
