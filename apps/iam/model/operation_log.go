package model

import (
	"gorm.io/gorm"
)

type OperationLogEntity struct {
	ID            uint           `gorm:"column:id;type:bigint;autoIncrement;primaryKey"`
	TenantID      uint           `gorm:"column:tenant_id;type:bigint;not null;default 0;comment: 租户ID"`
	UserID        uint           `gorm:"column:user_id;type:bigint;;comment: 操作人ID"`
	Username      string         `gorm:"column:username;type:varchar(32);;comment: 操作人账号"`
	Module        string         `gorm:"column:module;type:varchar(32);;comment: 操作模块"`
	Operation     string         `gorm:"column:operation;type:varchar(16);;comment: 操作类型: create/update/delete/query"`
	Method        string         `gorm:"column:method;type:varchar(16);;comment: 请求方法: GET/POST/PUT/DELETE等"`
	RequestID     string         `gorm:"column:request_id;type:varchar(64);;comment: 请求ID(用于追踪请求链路)"`
	RequestURL    string         `gorm:"column:request_url;type:varchar(512);;comment: 请求URL"`
	RequestParams string         `gorm:"column:request_params;type:text;;comment: 请求参数(JSON格式)"`
	ResponseResult string        `gorm:"column:response_result;type:text;;comment: 返回结果(JSON格式)"`
	IPAddress     string         `gorm:"column:ip_address;type:varchar(45);;comment: IP地址(支持IPv6)"`
	UserAgent     string         `gorm:"column:user_agent;type:varchar(512);;comment: 用户代理"`
	Status        string         `gorm:"column:status;type:varchar(16);;default success;comment: 操作状态: success-成功 failed-失败"`
	ErrorMsg      string         `gorm:"column:error_msg;type:varchar(1000);;comment: 错误信息"`
	ExecuteTime   int            `gorm:"column:execute_time;type:int;;comment: 执行时长(ms)"`
	CreatedAt     string         `gorm:"column:created_at;type:datetime(3);;comment: 创建时间"`
	UpdatedAt     string         `gorm:"column:updated_at;type:datetime(3);;comment: 更新时间"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3);;comment: 删除时间"`
	CreatedBy     uint           `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	UpdatedBy     uint           `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	DeletedBy     uint           `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
}

type OperationLogEntityList []OperationLogEntity

const TableNameOperationLog = "iam_operation_log"

func (OperationLogEntity) TableName() string {
	return TableNameOperationLog
}

func (l OperationLogEntityList) ToMap() map[uint]OperationLogEntity {
	m := make(map[uint]OperationLogEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
