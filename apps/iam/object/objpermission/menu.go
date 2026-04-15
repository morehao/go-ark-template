package objpermission

import "github.com/morehao/goark/apps/iam/iammodel"

type MenuBaseInfo struct {
	CacheType     iammodel.MenuCacheType  `json:"cacheType" form:"cacheType"`         // 缓存类型: enabled-启用 disabled-禁用
	TenantID      uint                    `json:"tenantID" form:"tenantID"`           // 所属租户ID
	ComponentPath string                  `json:"componentPath" form:"componentPath"` // 组件路径
	Icon          string                  `json:"icon" form:"icon"`                   // 菜单图标
	LinkType      iammodel.MenuLinkType   `json:"linkType" form:"linkType"`           // 链接类型: internal-内部链接 external-外部链接
	MenuCode      string                  `json:"menuCode" form:"menuCode"`           // 菜单编码
	MenuName      string                  `json:"menuName" form:"menuName"`           // 菜单名称
	MenuType      iammodel.MenuType       `json:"menuType" form:"menuType"`           // 菜单类型: directory-目录 menu-菜单 button-按钮
	ParentID      uint                    `json:"parentID" form:"parentID"`           // 父菜单ID
	Permission    string                  `json:"permission" form:"permission"`       // 权限标识: sys:user:add
	RoutePath     string                  `json:"routePath" form:"routePath"`         // 路由地址
	SortOrder     int32                   `json:"sortOrder" form:"sortOrder"`         // 排序
	Status        iammodel.MenuStatus     `json:"status" form:"status"`               // 状态: enabled-启用 disabled-停用
	Visibility    iammodel.MenuVisibility `json:"visibility" form:"visibility"`       // 可见性: visible-可见 hidden-隐藏
}
