package model

import (
	"time"

	"gorm.io/gorm"
)

type PersonGender string

const (
	PersonGenderMale    PersonGender = "male"
	PersonGenderFemale  PersonGender = "female"
	PersonGenderUnknown PersonGender = "unknown"
)

// PersonEntity 自然人用户表结构体
type PersonEntity struct {
	gorm.Model
	AvatarUrl    string       `gorm:"column:avatar_url;type:varchar(255);;default '';comment: 头像URL"`
	BirthDate    time.Time    `gorm:"column:birth_date;type:date;;default '';comment: 出生日期"`
	CreatedBy    uint         `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy    uint         `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Email        string       `gorm:"column:email;type:varchar(64);;default '';comment: 邮箱"`
	Gender       PersonGender `gorm:"column:gender;type:varchar(8);;default '';comment: 性别: male-男 female-女 unknown-未知"`
	Mobile       string       `gorm:"column:mobile;type:varchar(16);;default '';comment: 手机号"`
	PasswordHash string       `gorm:"column:password_hash;type:varchar(128);;default '';comment: 密码哈希(不存储盐值,盐值在应用层生成)"`
	RealName     string       `gorm:"column:real_name;type:varchar(32);not null;default '';comment: 真实姓名"`
	Remark       string       `gorm:"column:remark;type:varchar(500);;default '';comment: 备注"`
	UpdatedBy    uint         `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	Wechat       string       `gorm:"column:wechat;type:varchar(32);;default '';comment: 微信号"`
}

type PersonEntityList []PersonEntity

const TableNamePerson = "iam_person"

func (PersonEntity) TableName() string {
	return TableNamePerson
}

func (l PersonEntityList) ToMap() map[uint]PersonEntity {
	m := make(map[uint]PersonEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
