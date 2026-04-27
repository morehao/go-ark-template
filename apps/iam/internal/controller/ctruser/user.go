package ctruser

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtouser"
	"github.com/morehao/goark/apps/iam/internal/service/svcuser"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type UserCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	AssignDepartments(ctx *gin.Context)
	ListDepartments(ctx *gin.Context)
	AssignRoles(ctx *gin.Context)
	ListRoles(ctx *gin.Context)
	GetCurrentUserInfo(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	LoginHistory(ctx *gin.Context)
	Logout(ctx *gin.Context)
	PendingList(ctx *gin.Context)
	Approve(ctx *gin.Context)
	LoginByPassword(ctx *gin.Context)
	SelectTenant(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	Register(ctx *gin.Context)
	UnlockAccount(ctx *gin.Context)
	LoginLogCreate(ctx *gin.Context)
	LoginLogPageList(ctx *gin.Context)
}

type userCtr struct {
	userSvc     svcuser.UserSvc
	authSvc     svcuser.AuthSvc
	loginLogSvc svcuser.LoginLogSvc
}

var _ UserCtr = (*userCtr)(nil)

func NewUserCtr() UserCtr {
	return &userCtr{
		userSvc:     svcuser.NewUserSvc(),
		authSvc:     svcuser.NewAuthSvc(),
		loginLogSvc: svcuser.NewLoginLogSvc(),
	}
}

// Create 创建用户管理
// @Tags 用户管理
// @Summary 创建用户管理
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserCreateReq true "创建用户管理"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/create [post]
func (ctr *userCtr) Create(ctx *gin.Context) {
	var req dtouser.UserCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Delete 删除用户管理
// @Tags 用户管理
// @Summary 删除用户管理
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserDeleteReq true "删除用户管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/iam/user/delete [post]
func (ctr *userCtr) Delete(ctx *gin.Context) {
	var req dtouser.UserDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.userSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

// Update 修改用户管理
// @Tags 用户管理
// @Summary 修改用户管理
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserUpdateReq true "修改用户管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/iam/user/update [post]
func (ctr *userCtr) Update(ctx *gin.Context) {
	var req dtouser.UserUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

// Detail 用户管理详情
// @Tags 用户管理
// @Summary 用户管理详情
// @accept application/json
// @Produce application/json
// @Param req query dtouser.UserDetailReq true "用户管理详情"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/detail [get]
func (ctr *userCtr) Detail(ctx *gin.Context) {
	var req dtouser.UserDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// PageList 用户管理列表
// @Tags 用户管理
// @Summary 用户管理列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserPageListReq true "用户管理列表"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/pageList [post]
func (ctr *userCtr) PageList(ctx *gin.Context) {
	var req dtouser.UserPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// AssignDepartments 分配用户部门（全量替换）
// @Tags 用户管理
// @Summary 分配用户部门（全量替换）
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserDepartmentsAssignReq true "分配用户部门"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "分配成功"}"
// @Router /v1/iam/user/assignDepartments [post]
func (ctr *userCtr) AssignDepartments(ctx *gin.Context) {
	var req dtouser.UserDepartmentsAssignReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.AssignDepartments(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "分配成功")
}

// ListDepartments 获取用户部门列表
// @Tags 用户管理
// @Summary 获取用户部门列表
// @accept application/json
// @Produce application/json
// @Param req query dtouser.UserDepartmentsReq true "获取用户部门列表"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserDepartmentsResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/listDepartments [get]
func (ctr *userCtr) ListDepartments(ctx *gin.Context) {
	var req dtouser.UserDepartmentsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.ListDepartments(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// AssignRoles 分配用户角色
// @Tags 用户管理
// @Summary 分配用户角色
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserAssignRolesReq true "分配用户角色"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "分配成功"}"
// @Router /v1/iam/user/assignRoles [post]
func (ctr *userCtr) AssignRoles(ctx *gin.Context) {
	var req dtouser.UserAssignRolesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.AssignRoles(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "分配成功")
}

// ListRoles 查询用户角色列表
// @Tags 用户管理
// @Summary 查询用户角色列表
// @accept application/json
// @Produce application/json
// @Param req query dtouser.UserRolesReq true "查询用户角色列表"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserRolesResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/listRoles [get]
func (ctr *userCtr) ListRoles(ctx *gin.Context) {
	var req dtouser.UserRolesReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.ListRoles(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// GetCurrentUserInfo 获取当前用户信息
// @Tags 用户管理
// @Summary 获取当前用户信息
// @accept application/json
// @Produce application/json
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserInfoResp}
// @Router /v1/iam/user/getCurrentUserInfo [get]
func (ctr *userCtr) GetCurrentUserInfo(ctx *gin.Context) {
	res, err := ctr.userSvc.GetCurrentUserInfo(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// UpdateProfile 更新个人资料
// @Tags 用户管理
// @Summary 更新个人资料
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UpdateProfileReq true "更新个人资料"
// @Success 200 {object} gincontext.DtoRender{data=string}
// @Router /v1/iam/user/updateProfile [post]
func (ctr *userCtr) UpdateProfile(ctx *gin.Context) {
	var req dtouser.UpdateProfileReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.UpdateProfile(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "更新成功")
}

// ChangePassword 修改密码
// @Tags 用户管理
// @Summary 修改密码
// @accept application/json
// @Produce application/json
// @Param req body dtouser.ChangePasswordReq true "修改密码"
// @Success 200 {object} gincontext.DtoRender{data=string}
// @Router /v1/iam/user/changePassword [post]
func (ctr *userCtr) ChangePassword(ctx *gin.Context) {
	var req dtouser.ChangePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.ChangePassword(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "密码修改成功")
}

// LoginHistory 查看登录历史
// @Tags 用户管理
// @Summary 查看登录历史
// @accept application/json
// @Produce application/json
// @Param req query dtouser.LoginHistoryReq true "查看登录历史"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.LoginHistoryResp}
// @Router /v1/iam/user/loginHistory [get]
func (ctr *userCtr) LoginHistory(ctx *gin.Context) {
	var req dtouser.LoginHistoryReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.LoginHistory(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Logout 登出
// @Tags 用户管理
// @Summary 登出
// @accept application/json
// @Produce application/json
// @Success 200 {object} gincontext.DtoRender{data=string}
// @Router /v1/iam/user/logout [post]
func (ctr *userCtr) Logout(ctx *gin.Context) {
	if err := ctr.userSvc.Logout(ctx); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "登出成功")
}

// PendingList 待审批用户列表
// @Tags 用户管理
// @Summary 待审批用户列表
// @accept application/json
// @Produce application/json
// @Param req query dtouser.PendingListReq true "待审批用户列表"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.PendingListResp}
// @Router /v1/iam/user/pendingList [get]
func (ctr *userCtr) PendingList(ctx *gin.Context) {
	var req dtouser.PendingListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.userSvc.PendingList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Approve 审批用户
// @Tags 用户管理
// @Summary 审批用户
// @accept application/json
// @Produce application/json
// @Param req body dtouser.ApproveReq true "审批用户"
// @Success 200 {object} gincontext.DtoRender{data=string}
// @Router /v1/iam/user/approve [post]
func (ctr *userCtr) Approve(ctx *gin.Context) {
	var req dtouser.ApproveReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.Approve(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "审批成功")
}

// LoginByPassword 密码登录
// @Tags 用户管理
// @Summary 密码登录
// @accept application/json
// @Produce application/json
// @Param req body dtouser.LoginByPasswordReq true "密码登录请求"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.LoginByPasswordResp}
// @Router /v1/iam/user/loginByPassword [post]
func (ctr *userCtr) LoginByPassword(ctx *gin.Context) {
	var req dtouser.LoginByPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.authSvc.LoginByPassword(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// SelectTenant 选择租户
// @Tags 用户管理
// @Summary 选择租户
// @accept application/json
// @Produce application/json
// @Param req body dtouser.SelectTenantReq true "选择租户请求"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.SelectTenantResp}
// @Router /v1/iam/user/selectTenant [post]
func (ctr *userCtr) SelectTenant(ctx *gin.Context) {
	var req dtouser.SelectTenantReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.authSvc.SelectTenant(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// RefreshToken 刷新令牌
// @Tags 用户管理
// @Summary 刷新令牌
// @accept application/json
// @Produce application/json
// @Param req body dtouser.RefreshTokenReq true "刷新令牌请求"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.RefreshTokenResp}
// @Router /v1/iam/user/refreshToken [post]
func (ctr *userCtr) RefreshToken(ctx *gin.Context) {
	var req dtouser.RefreshTokenReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.authSvc.RefreshToken(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Register 用户注册
// @Tags 用户管理
// @Summary 用户注册
// @accept application/json
// @Produce application/json
// @Param req body dtouser.RegisterReq true "用户注册请求"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.RegisterResp}
// @Router /v1/iam/user/register [post]
func (ctr *userCtr) Register(ctx *gin.Context) {
	var req dtouser.RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.authSvc.Register(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// UnlockAccount 自助解锁账户
// @Tags 用户管理
// @Summary 自助解锁账户
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UnlockAccountReq true "自助解锁账户"
// @Success 200 {object} gincontext.DtoRender{data=string}
// @Router /v1/iam/user/unlockAccount [post]
func (ctr *userCtr) UnlockAccount(ctx *gin.Context) {
	var req dtouser.UnlockAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.authSvc.UnlockAccount(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "解锁成功")
}

// LoginLogCreate 创建登录日志
// @Tags 用户管理
// @Summary 创建登录日志
// @accept application/json
// @Produce application/json
// @Param req body dtouser.LoginLogCreateReq true "创建登录日志"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.LoginLogCreateResp}
// @Router /v1/iam/user/loginLog/create [post]
func (ctr *userCtr) LoginLogCreate(ctx *gin.Context) {
	var req dtouser.LoginLogCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.loginLogSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// LoginLogPageList 登录日志列表
// @Tags 用户管理
// @Summary 登录日志列表
// @accept application/json
// @Produce application/json
// @Param req query dtouser.LoginLogPageListReq true "登录日志列表"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.LoginLogPageListResp}
// @Router /v1/iam/user/loginLog/pageList [get]
func (ctr *userCtr) LoginLogPageList(ctx *gin.Context) {
	var req dtouser.LoginLogPageListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.loginLogSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
