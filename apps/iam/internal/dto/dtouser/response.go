package dtouser

import (
	"github.com/morehao/goark/iam/object/objuser"
	"github.com/morehao/golib/biz/gobject"
)

type UserCreateResp struct {
	UserID   uint   `json:"userID"`             // 用户UserID
	PersonID uint   `json:"personID"`           // 自然人ID
	Password string `json:"password,omitempty"` // 仅创建成功时返回
}

type UserDetailResp struct {
	UserID uint `json:"userID" validate:"required"` // 数据自增 ID
	objuser.UserBaseInfo
	gobject.OperatorBaseInfo
}

type UserPageListItem struct {
	UserID uint `json:"userID" validate:"required"` // 数据自增 ID
	objuser.UserBaseInfo
	gobject.OperatorBaseInfo
}

type UserPageListResp struct {
	List  []UserPageListItem `json:"list"`  // 数据列表
	Total int64              `json:"total"` // 数据总条数
}

type UserDepartmentItem struct {
	DepartmentID   uint   `json:"departmentID"`   // 部门ID
	DepartmentName string `json:"departmentName"` // 部门名称
	DeptType       string `json:"deptType"`       // 部门类型
}

type UserDepartmentsResp struct {
	List []UserDepartmentItem `json:"list"` // 用户部门列表
}

type UserRoleItem struct {
	RoleID   uint   `json:"roleID"`   // 角色ID
	RoleName string `json:"roleName"` // 角色名称
	RoleCode string `json:"roleCode"` // 角色编码
	RoleType string `json:"roleType"` // 角色类型
}

type UserRolesResp struct {
	List []UserRoleItem `json:"list"` // 角色列表
}

type UserInfoResp struct {
	UserID     uint     `json:"userID"`     // 用户ID
	Username   string   `json:"username"`   // 用户名
	PersonID   uint     `json:"personID"`   // 自然人ID
	Email      string   `json:"email"`      // 邮箱
	Phone      string   `json:"phone"`      // 手机号
	Avatar     string   `json:"avatar"`     // 头像
	Nickname   string   `json:"nickname"`   // 昵称
	Status     string   `json:"status"`     // 状态
	UserType   string   `json:"userType"`   // 用户类型
	TenantID   uint     `json:"tenantID"`   // 租户ID
	TenantName string   `json:"tenantName"` // 租户名称
	OrgID      uint     `json:"orgID"`      // 组织ID
	OrgName    string   `json:"orgName"`    // 组织名称
	RoleIDs    []uint   `json:"roleIDs"`    // 角色ID列表
	RoleNames  []string `json:"roleNames"`  // 角色名称列表
	DeptIDs    []uint   `json:"deptIDs"`    // 部门ID列表
	DeptNames  []string `json:"deptNames"`  // 部门名称列表
}

type LoginHistoryItem struct {
	ID           uint   `json:"id"`            // 日志ID
	LoginType    string `json:"loginType"`     // 登录类型
	LoginStatus  string `json:"loginStatus"`   // 登录状态
	LoginMessage string `json:"loginMessage"`  // 登录消息
	IPAddress    string `json:"ipAddress"`     // IP地址
	Location     string `json:"location"`      // 位置
	Browser      string `json:"browser"`       // 浏览器
	OS           string `json:"os"`            // 操作系统
	CreatedAt    string `json:"createdAt"`     // 创建时间
}

type LoginHistoryResp struct {
	List  []LoginHistoryItem `json:"list"`  // 数据列表
	Total int64              `json:"total"` // 数据总条数
}
