package model

import (
	"gorm.io/gorm"
)

type LoginLogEntity struct {
	ID           uint           `gorm:"column:id;type:bigint;autoIncrement;primaryKey"`
	TenantID     uint           `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 租户ID"`
	UserID       uint           `gorm:"column:user_id;type:bigint;;comment: 用户ID"`
	Username     string         `gorm:"column:username;type:varchar(32);;comment: 用户名"`
	LoginType    string         `gorm:"column:login_type;type:varchar(16);;comment: 登录类型: password/sms/wechat"`
	LoginStatus  string         `gorm:"column:login_status;type:varchar(16);;comment: 登录状态: success-成功 failed-失败"`
	LoginMessage string         `gorm:"column:login_message;type:varchar(128);;comment: 登录消息"`
	IPAddress    string         `gorm:"column:ip_address;type:varchar(45);;comment: IP地址(支持IPv6)"`
	Location     string         `gorm:"column:location;type:varchar(128);;comment: 登录地点"`
	Browser      string         `gorm:"column:browser;type:varchar(64);;comment: 浏览器"`
	OS           string         `gorm:"column:os;type:varchar(64);;comment: 操作系统"`
	CreatedAt    string         `gorm:"column:created_at;type:datetime(3);;comment: 创建时间"`
	UpdatedAt    string         `gorm:"column:updated_at;type:datetime(3);;comment: 更新时间"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);;comment: 删除时间"`
	CreatedBy    uint           `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	UpdatedBy    uint           `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	DeletedBy    uint           `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
}

type LoginLogEntityList []LoginLogEntity

const TableNameLoginLog = "iam_login_log"

func (LoginLogEntity) TableName() string {
	return TableNameLoginLog
}

func (l LoginLogEntityList) ToMap() map[uint]LoginLogEntity {
	m := make(map[uint]LoginLogEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
