package objpermission

type MenuBaseInfo struct {
	// CacheType 缓存类型: enabled-启用 disabled-禁用
	CacheType string `json:"cacheType" form:"cacheType"`
	// TenantID 所属租户ID
	TenantID uint `json:"tenantID" form:"tenantID"`
	// ComponentPath 组件路径
	ComponentPath string `json:"componentPath" form:"componentPath"`
	// Icon 菜单图标
	Icon string `json:"icon" form:"icon"`
	// LinkType 链接类型: internal-内部链接 external-外部链接
	LinkType string `json:"linkType" form:"linkType"`
	// MenuCode 菜单编码
	MenuCode string `json:"menuCode" form:"menuCode"`
	// MenuName 菜单名称
	MenuName string `json:"menuName" form:"menuName"`
	// MenuType 菜单类型: directory-目录 menu-菜单 button-按钮
	MenuType string `json:"menuType" form:"menuType"`
	// ParentID 父菜单ID
	ParentID uint `json:"parentID" form:"parentID"`
	// Permission 权限标识: sys:user:add
	Permission string `json:"permission" form:"permission"`
	// RoutePath 路由地址
	RoutePath string `json:"routePath" form:"routePath"`
	// SortOrder 排序
	SortOrder int32 `json:"sortOrder" form:"sortOrder"`
	// Status 状态: enabled-启用 disabled-停用
	Status string `json:"status" form:"status"`
	// Visibility 可见性: visible-可见 hidden-隐藏
	Visibility string `json:"visibility" form:"visibility"`
}
