package ctrpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/apps/iam/internal/service/svcpermission"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type AuthorizationCtr interface {
	AssignRoles(ctx *gin.Context)
	RemoveRoles(ctx *gin.Context)
	ListUserRoles(ctx *gin.Context)
	GetUserPermissions(ctx *gin.Context)
}

type authorizationCtr struct {
	authorizationSvc svcpermission.AuthorizationSvc
}

var _ AuthorizationCtr = (*authorizationCtr)(nil)

func NewAuthorizationCtr() AuthorizationCtr {
	return &authorizationCtr{
		authorizationSvc: svcpermission.NewAuthorizationSvc(),
	}
}

// AssignRoles 为用户分配角色
// @Tags 权限管理
// @Summary 为用户分配角色
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.UserAssignRolesReq true "为用户分配角色"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "分配成功"}"
// @Router /v1/iam/user/assignRoles [post]
func (ctr *authorizationCtr) AssignRoles(ctx *gin.Context) {
	var req dtopermission.UserAssignRolesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.authorizationSvc.AssignRolesToUser(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "分配成功")
	}
}

// RemoveRoles 移除用户角色
// @Tags 权限管理
// @Summary 移除用户角色
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.UserRemoveRolesReq true "移除用户角色"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "移除成功"}"
// @Router /v1/iam/user/removeRoles [post]
func (ctr *authorizationCtr) RemoveRoles(ctx *gin.Context) {
	var req dtopermission.UserRemoveRolesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.authorizationSvc.RemoveRolesFromUser(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "移除成功")
	}
}

// ListUserRoles 获取用户角色列表
// @Tags 权限管理
// @Summary 获取用户角色列表
// @accept application/json
// @Produce application/json
// @Param req query dtopermission.UserRoleListReq true "获取用户角色列表"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.UserRoleListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/user/listRoles [get]
func (ctr *authorizationCtr) ListUserRoles(ctx *gin.Context) {
	var req dtopermission.UserRoleListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.authorizationSvc.ListUserRoles(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// GetUserPermissions 获取当前用户权限
// @Tags 权限管理
// @Summary 获取当前用户权限
// @accept application/json
// @Produce application/json
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.UserPermissionsResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/permission/userPermissions [get]
func (ctr *authorizationCtr) GetUserPermissions(ctx *gin.Context) {
	res, err := ctr.authorizationSvc.GetUserPermissions(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
