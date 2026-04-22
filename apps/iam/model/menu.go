package model

import (
	"strings"

	"gorm.io/gorm"
)

type MenuCacheType string

const (
	MenuCacheTypeEnabled  MenuCacheType = "enabled"
	MenuCacheTypeDisabled MenuCacheType = "disabled"
)

type MenuLinkType string

const (
	MenuLinkTypeInternal MenuLinkType = "internal"
	MenuLinkTypeExternal MenuLinkType = "external"
)

type MenuType string

const (
	MenuTypeDirectory MenuType = "directory"
	MenuTypeMenu      MenuType = "menu"
	MenuTypeButton    MenuType = "button"
)

type MenuStatus string

const (
	MenuStatusEnabled  MenuStatus = "enabled"
	MenuStatusDisabled MenuStatus = "disabled"
)

type MenuVisibility string

const (
	MenuVisibilityVisible MenuVisibility = "visible"
	MenuVisibilityHidden  MenuVisibility = "hidden"
)

type MenuAccessPolicy int

const (
	AccessPolicyPublic      MenuAccessPolicy = 1 << iota // 1   = 0001
	AccessPolicyAuthorized                               // 2   = 0010
	AccessPolicyOrgAdmin                                 // 4   = 0100
	AccessPolicyTenantAdmin                              // 8   = 1000
)

type MenuAccessPolicyString string

const (
	MenuAccessPolicyPublicStr      MenuAccessPolicyString = "public"
	MenuAccessPolicyAuthorizedStr  MenuAccessPolicyString = "authorized"
	MenuAccessPolicyOrgAdminStr    MenuAccessPolicyString = "org_admin"
	MenuAccessPolicyTenantAdminStr MenuAccessPolicyString = "tenant_admin"
)

func (m MenuAccessPolicy) HasPolicy(policy MenuAccessPolicy) bool {
	return int(m)&int(policy) != 0
}

func (m *MenuAccessPolicy) AddPolicy(policy MenuAccessPolicy) {
	*m |= MenuAccessPolicy(policy)
}

func (m *MenuAccessPolicy) RemovePolicy(policy MenuAccessPolicy) {
	*m &= ^MenuAccessPolicy(policy)
}

func (m *MenuAccessPolicy) SetPolicy(policy MenuAccessPolicy) {
	*m = policy
}

func (m MenuAccessPolicy) IsEmpty() bool {
	return m == 0
}

func (m MenuAccessPolicy) ToStrings() []MenuAccessPolicyString {
	var policies []MenuAccessPolicyString
	if m.HasPolicy(AccessPolicyPublic) {
		policies = append(policies, MenuAccessPolicyPublicStr)
	}
	if m.HasPolicy(AccessPolicyAuthorized) {
		policies = append(policies, MenuAccessPolicyAuthorizedStr)
	}
	if m.HasPolicy(AccessPolicyOrgAdmin) {
		policies = append(policies, MenuAccessPolicyOrgAdminStr)
	}
	if m.HasPolicy(AccessPolicyTenantAdmin) {
		policies = append(policies, MenuAccessPolicyTenantAdminStr)
	}
	return policies
}

func AccessPoliciesToMask(policies []MenuAccessPolicyString) MenuAccessPolicy {
	var mask MenuAccessPolicy
	for _, s := range policies {
		switch MenuAccessPolicyString(s) {
		case MenuAccessPolicyPublicStr:
			mask.AddPolicy(AccessPolicyPublic)
		case MenuAccessPolicyAuthorizedStr:
			mask.AddPolicy(AccessPolicyAuthorized)
		case MenuAccessPolicyOrgAdminStr:
			mask.AddPolicy(AccessPolicyOrgAdmin)
		case MenuAccessPolicyTenantAdminStr:
			mask.AddPolicy(AccessPolicyTenantAdmin)
		}
	}
	return mask
}

func (m MenuAccessPolicy) String() string {
	var policies []string
	if m.HasPolicy(AccessPolicyPublic) {
		policies = append(policies, "public")
	}
	if m.HasPolicy(AccessPolicyAuthorized) {
		policies = append(policies, "authorized")
	}
	if m.HasPolicy(AccessPolicyOrgAdmin) {
		policies = append(policies, "org_admin")
	}
	if m.HasPolicy(AccessPolicyTenantAdmin) {
		policies = append(policies, "tenant_admin")
	}
	if len(policies) == 0 {
		return "none"
	}
	return strings.Join(policies, "|")
}

// MenuEntity 菜单管理表结构体
type MenuEntity struct {
	gorm.Model
	CacheType     MenuCacheType    `gorm:"column:cache_type;type:varchar(16);;default disabled;comment: 缓存类型: enabled-启用 disabled-禁用"`
	TenantID      uint             `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 所属租户ID"`
	ComponentPath string           `gorm:"column:component_path;type:varchar(255);;default '';comment: 组件路径"`
	CreatedBy     uint             `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy     uint             `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Icon          string           `gorm:"column:icon;type:varchar(64);;default '';comment: 菜单图标"`
	LinkType      MenuLinkType     `gorm:"column:link_type;type:varchar(16);;default internal;comment: 链接类型: internal-内部链接 external-外部链接"`
	MenuCode      string           `gorm:"column:menu_code;type:varchar(32);not null;default '';comment: 菜单编码"`
	MenuName      string           `gorm:"column:menu_name;type:varchar(64);not null;default '';comment: 菜单名称"`
	MenuType      MenuType         `gorm:"column:menu_type;type:varchar(16);;default directory;comment: 菜单类型: directory-目录 menu-菜单 button-按钮"`
	ParentID      uint             `gorm:"column:parent_id;type:bigint;;default 0;comment: 父菜单ID"`
	Permission    string           `gorm:"column:permission;type:varchar(64);;default '';comment: 权限标识: sys:user:add"`
	RoutePath     string           `gorm:"column:route_path;type:varchar(255);;default '';comment: 路由地址"`
	SortOrder     int32            `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	Status        MenuStatus       `gorm:"column:status;type:varchar(16);;default enabled;comment: 状态: enabled-启用 disabled-停用"`
	UpdatedBy     uint             `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	Visibility    MenuVisibility   `gorm:"column:visibility;type:varchar(16);;default visible;comment: 是否显示: visible-显示 hidden-隐藏"`
	AccessPolicy  MenuAccessPolicy `gorm:"column:access_policy;type:int;;default 1;comment: 访问策略位掩码: 1-全部人可见 2-需授权 4-组织管理员 8-租户管理员"`
}

type MenuEntityList []MenuEntity

const TableNameMenu = "iam_menu"

func (MenuEntity) TableName() string {
	return TableNameMenu
}

func (l MenuEntityList) ToMap() map[uint]MenuEntity {
	m := make(map[uint]MenuEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}
