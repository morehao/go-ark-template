package model

import (
	"gorm.io/gorm"
)

// OrganizationConfigEntity 组织配置表结构体
type OrganizationConfigEntity struct {
	gorm.Model
	ConfigGroup string `gorm:"column:config_group;type:varchar(32);;default general;comment: 配置分组: general-通用/auth-认证/theme-主题等"`
	ConfigKey   string `gorm:"column:config_key;type:varchar(100);not null;default '';comment: 配置键"`
	ValueType  string `gorm:"column:value_type;type:varchar(32);;default string;comment: 配置类型: string/json/boolean/number"`
	ConfigValue string `gorm:"column:config_value;type:text;;default '';comment: 配置值(支持JSON)"`
	Description string `gorm:"column:description;type:varchar(255);not null;default '';comment: 配置说明"`
	Sequence   int32  `gorm:"column:sequence;type:int;;default 0;comment: 排序"`
	OrgID       uint   `gorm:"column:org_id;type:bigint;not null;default 0;comment: 组织ID"`
}

type OrganizationConfigEntityList []OrganizationConfigEntity

const TableNameOrganizationConfig = "iam_organization_config"

func (OrganizationConfigEntity) TableName() string {
	return TableNameOrganizationConfig
}

func (l OrganizationConfigEntityList) ToMap() map[uint]OrganizationConfigEntity {
	m := make(map[uint]OrganizationConfigEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

const (
	OrgConfigGroupAuth    = "auth"
	OrgConfigGroupGeneral = "general"
	OrgConfigGroupTheme   = "theme"
)

const (
	OrgConfigTypeString  = "string"
	OrgConfigTypeBoolean = "boolean"
	OrgConfigTypeNumber  = "number"
	OrgConfigTypeJSON    = "json"
)

const (
	OrgConfigKeyRegisterEnabled         = "auth.register.enabled"
	OrgConfigKeyRegisterRequireApproval = "auth.register.requireApproval"
	OrgConfigKeyRegisterIdentityType    = "auth.register.identityType"
	OrgConfigKeyPasswordMinLength       = "auth.password.minLength"
	OrgConfigKeyPasswordRequireUppercase = "auth.password.requireUppercase"
	OrgConfigKeyPasswordRequireLowercase = "auth.password.requireLowercase"
	OrgConfigKeyPasswordRequireNumber   = "auth.password.requireNumber"
	OrgConfigKeyPasswordRequireSpecial   = "auth.password.requireSpecial"
	OrgConfigKeyLoginMaxFailCount       = "auth.login.maxFailCount"
	OrgConfigKeyLoginLockDuration       = "auth.login.lockDuration"
)

type RegisterIdentityType string

const (
	RegisterIdentityTypeMobile RegisterIdentityType = "mobile"
	RegisterIdentityTypeEmail  RegisterIdentityType = "email"
	RegisterIdentityTypeBoth   RegisterIdentityType = "both"
)

type OrgConfigMeta struct {
	Group        string            `json:"group"`
	Key          string            `json:"key"`
	Type         string            `json:"type"`
	DefaultValue string            `json:"defaultValue"`
	Description  string            `json:"description"`
	Options      []OrgConfigOption `json:"options,omitempty"`
}

type OrgConfigOption struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

var OrgConfigMetaList = []OrgConfigMeta{
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyRegisterEnabled,
		Type:         OrgConfigTypeBoolean,
		DefaultValue: "true",
		Description:  "是否开放注册",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyRegisterRequireApproval,
		Type:         OrgConfigTypeBoolean,
		DefaultValue: "false",
		Description:  "注册是否需要审核",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyRegisterIdentityType,
		Type:         OrgConfigTypeString,
		DefaultValue: string(RegisterIdentityTypeEmail),
		Description:  "注册身份认证类型",
		Options: []OrgConfigOption{
			{Value: string(RegisterIdentityTypeMobile), Description: "手机号注册"},
			{Value: string(RegisterIdentityTypeEmail), Description: "邮箱注册"},
			{Value: string(RegisterIdentityTypeBoth), Description: "手机号或邮箱注册"},
		},
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyPasswordMinLength,
		Type:         OrgConfigTypeNumber,
		DefaultValue: "8",
		Description:  "密码最小长度",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyPasswordRequireUppercase,
		Type:         OrgConfigTypeBoolean,
		DefaultValue: "true",
		Description:  "密码必须包含大写字母",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyPasswordRequireLowercase,
		Type:         OrgConfigTypeBoolean,
		DefaultValue: "true",
		Description:  "密码必须包含小写字母",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyPasswordRequireNumber,
		Type:         OrgConfigTypeBoolean,
		DefaultValue: "true",
		Description:  "密码必须包含数字",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyPasswordRequireSpecial,
		Type:         OrgConfigTypeBoolean,
		DefaultValue: "false",
		Description:  "密码必须包含特殊字符",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyLoginMaxFailCount,
		Type:         OrgConfigTypeNumber,
		DefaultValue: "5",
		Description:  "登录失败最大次数",
	},
	{
		Group:        OrgConfigGroupAuth,
		Key:          OrgConfigKeyLoginLockDuration,
		Type:         OrgConfigTypeNumber,
		DefaultValue: "300",
		Description:  "登录锁定时长(秒)",
	},
}

func (m *OrgConfigMeta) ValidateValue(value string) bool {
	if len(m.Options) == 0 {
		return true
	}
	for _, opt := range m.Options {
		if opt.Value == value {
			return true
		}
	}
	return false
}

func GetOrgConfigMetaByKey(key string) *OrgConfigMeta {
	for i := range OrgConfigMetaList {
		if OrgConfigMetaList[i].Key == key {
			return &OrgConfigMetaList[i]
		}
	}
	return nil
}
