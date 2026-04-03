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
	AssignDepartment(ctx *gin.Context)
	RemoveDepartment(ctx *gin.Context)
	ListDepartments(ctx *gin.Context)
	AssignRoles(ctx *gin.Context)
	RemoveRoles(ctx *gin.Context)
	ListRoles(ctx *gin.Context)
}

type userCtr struct {
	userSvc svcuser.UserSvc
}

var _ UserCtr = (*userCtr)(nil)

func NewUserCtr() UserCtr {
	return &userCtr{
		userSvc: svcuser.NewUserSvc(),
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
// @Param req query dtouser.UserPageListReq true "用户管理列表"
// @Success 200 {object} gincontext.DtoRender{data=dtouser.UserPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/pageList [post]
func (ctr *userCtr) PageList(ctx *gin.Context) {
	var req dtouser.UserPageListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
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

// AssignDepartment 分配用户部门
// @Tags 用户管理
// @Summary 分配用户部门
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserDepartmentAssignReq true "分配用户部门"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/assignDepartment [post]
func (ctr *userCtr) AssignDepartment(ctx *gin.Context) {
	var req dtouser.UserDepartmentAssignReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.AssignDepartment(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "分配成功")
}

// RemoveDepartment 移除用户部门
// @Tags 用户管理
// @Summary 移除用户部门
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserDepartmentRemoveReq true "移除用户部门"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/removeDepartment [post]
func (ctr *userCtr) RemoveDepartment(ctx *gin.Context) {
	var req dtouser.UserDepartmentRemoveReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.RemoveDepartment(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "移除成功")
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

// RemoveRoles 移除用户角色
// @Tags 用户管理
// @Summary 移除用户角色
// @accept application/json
// @Produce application/json
// @Param req body dtouser.UserRemoveRolesReq true "移除用户角色"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "移除成功"}"
// @Router /v1/iam/user/removeRoles [post]
func (ctr *userCtr) RemoveRoles(ctx *gin.Context) {
	var req dtouser.UserRemoveRolesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.userSvc.RemoveRoles(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "移除成功")
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
