package ctrpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/apps/iam/internal/service/svcpermission"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type RoleCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	AssignMenus(ctx *gin.Context)
	ListMenus(ctx *gin.Context)
}

type roleCtr struct {
	roleSvc svcpermission.RoleSvc
}

var _ RoleCtr = (*roleCtr)(nil)

func NewRoleCtr() RoleCtr {
	return &roleCtr{
		roleSvc: svcpermission.NewRoleSvc(),
	}
}

// Create 创建角色管理
// @Tags 角色管理
// @Summary 创建角色管理
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.RoleCreateReq true "创建角色管理"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.RoleCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/role/create [post]
func (ctr *roleCtr) Create(ctx *gin.Context) {
	var req dtopermission.RoleCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.roleSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Delete 删除角色管理
// @Tags 角色管理
// @Summary 删除角色管理
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.RoleDeleteReq true "删除角色管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/iam/role/delete [post]
func (ctr *roleCtr) Delete(ctx *gin.Context) {
	var req dtopermission.RoleDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.roleSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

// Update 修改角色管理
// @Tags 角色管理
// @Summary 修改角色管理
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.RoleUpdateReq true "修改角色管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/iam/role/update [post]
func (ctr *roleCtr) Update(ctx *gin.Context) {
	var req dtopermission.RoleUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.roleSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

// Detail 角色管理详情
// @Tags 角色管理
// @Summary 角色管理详情
// @accept application/json
// @Produce application/json
// @Param req query dtopermission.RoleDetailReq true "角色管理详情"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.RoleDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/role/detail [get]
func (ctr *roleCtr) Detail(ctx *gin.Context) {
	var req dtopermission.RoleDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.roleSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// PageList 角色管理列表
// @Tags 角色管理
// @Summary 角色管理列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.RolePageListReq true "角色管理列表"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.RolePageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/role/pageList [post]
func (ctr *roleCtr) PageList(ctx *gin.Context) {
	var req dtopermission.RolePageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.roleSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// AssignMenus 角色分配菜单
// @Tags 角色管理
// @Summary 角色分配菜单
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.RoleAssignMenusReq true "角色分配菜单"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "分配成功"}"
// @Router /v1/iam/role/assignMenus [post]
func (ctr *roleCtr) AssignMenus(ctx *gin.Context) {
	var req dtopermission.RoleAssignMenusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.roleSvc.AssignMenus(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "分配成功")
}

// ListMenus 查询角色菜单列表
// @Tags 角色管理
// @Summary 查询角色菜单列表
// @accept application/json
// @Produce application/json
// @Param req query dtopermission.RoleListMenusReq true "查询角色菜单"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.RoleMenuListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/role/listMenus [get]
func (ctr *roleCtr) ListMenus(ctx *gin.Context) {
	var req dtopermission.RoleListMenusReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.roleSvc.ListMenus(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
