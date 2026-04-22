package objpermission

import "github.com/morehao/goark/apps/iam/model"

type RoleBaseInfo struct {
	TenantID    uint                   `json:"tenantID" form:"tenantID"`       // 所属租户ID
	DataScope   model.RoleDataScope `json:"dataScope" form:"dataScope"`     // 数据权限范围: all-全部 dept_and_sub-本部门及以下 dept-本部门 self-仅本人 custom-自定义
	Description string                 `json:"description" form:"description"` // 角色描述
	RoleCode    string                 `json:"roleCode" form:"roleCode"`       // 角色编码
	RoleName    string                 `json:"roleName" form:"roleName"`       // 角色名称
	RoleType    model.RoleType      `json:"roleType" form:"roleType"`       // 角色类型: custom-自定义 system-系统内置
	Sequence   int32                  `json:"sequence" form:"sequence"`     // 排序
	Status      model.RoleStatus    `json:"status" form:"status"`           // 状态: enabled-启用 disabled-停用
}
