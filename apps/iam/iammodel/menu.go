package iammodel

import (
	"gorm.io/gorm"
)

// MenuEntity 菜单管理表结构体
type MenuEntity struct {
	gorm.Model
	CacheType     string `gorm:"column:cache_type;type:varchar(16);;default disabled;comment: 缓存类型: enabled-启用 disabled-禁用"`
	TenantID      uint   `gorm:"column:tenant_id;type:bigint;not null;default '';comment: 所属租户ID"`
	ComponentPath string `gorm:"column:component_path;type:varchar(255);;default '';comment: 组件路径"`
	CreatedBy     uint   `gorm:"column:created_by;type:bigint;not null;default 0;comment: 创建人ID"`
	DeletedBy     uint   `gorm:"column:deleted_by;type:bigint;not null;default 0;comment: 删除人ID"`
	Icon          string `gorm:"column:icon;type:varchar(64);;default '';comment: 菜单图标"`
	LinkType      string `gorm:"column:link_type;type:varchar(16);;default internal;comment: 链接类型: internal-内部链接 external-外部链接"`
	MenuCode      string `gorm:"column:menu_code;type:varchar(32);not null;default '';comment: 菜单编码"`
	MenuName      string `gorm:"column:menu_name;type:varchar(64);not null;default '';comment: 菜单名称"`
	MenuType      string `gorm:"column:menu_type;type:varchar(16);;default directory;comment: 菜单类型: directory-目录 menu-菜单 button-按钮"`
	ParentID      uint   `gorm:"column:parent_id;type:bigint;;default 0;comment: 父菜单ID"`
	Permission    string `gorm:"column:permission;type:varchar(64);;default '';comment: 权限标识: sys:user:add"`
	RoutePath     string `gorm:"column:route_path;type:varchar(255);;default '';comment: 路由地址"`
	SortOrder     int32  `gorm:"column:sort_order;type:int;;default 0;comment: 排序"`
	Status        string `gorm:"column:status;type:varchar(16);;default active;comment: 状态: active-正常 inactive-停用"`
	UpdatedBy     uint   `gorm:"column:updated_by;type:bigint;not null;default 0;comment: 更新人ID"`
	Visibility    string `gorm:"column:visibility;type:varchar(16);;default visible;comment: 可见性: visible-可见 hidden-隐藏"`
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
