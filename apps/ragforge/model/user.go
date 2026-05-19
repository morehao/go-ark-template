package model

import (
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin      UserRole = "admin"
	UserRoleUser       UserRole = "user"
	UserRoleViewer     UserRole = "viewer"
)

type UserEntity struct {
	gorm.Model
	TenantID     uint     `gorm:"column:tenant_id;type:bigint unsigned;not null;default:0;index;comment:租户id"`
	Username     string   `gorm:"column:username;type:varchar(255);not null;default:'';comment:用户名"`
	Email        string   `gorm:"column:email;type:varchar(255);not null;default:'';comment:邮箱"`
	PasswordHash string   `gorm:"column:password_hash;type:varchar(500);not null;default:'';comment:密码哈希"`
	Role         UserRole `gorm:"column:role;type:varchar(50);not null;default:'user';comment:角色"`
}

type UserEntityList []UserEntity

func (UserEntity) TableName() string {
	return TableNameUser
}

