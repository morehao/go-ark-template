package ctrauth

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/apps/iam/internal/service/svcauth"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type AuthCtr interface {
	Login(ctx *gin.Context)
	SelectTenant(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authCtr struct {
	authSvc svcauth.AuthSvc
}

var _ AuthCtr = (*authCtr)(nil)

func NewAuthCtr() AuthCtr {
	return &authCtr{
		authSvc: svcauth.NewAuthSvc(),
	}
}

// Login 用户登录
// @Tags 认证管理
// @Summary 用户登录
// @accept application/json
// @Produce application/json
// @Param req body dtoauth.LoginReq true "登录请求"
// @Success 200 {object} gincontext.DtoRender{data=dtoauth.LoginResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/auth/login [post]
func (ctr *authCtr) Login(ctx *gin.Context) {
	var req dtoauth.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.authSvc.Login(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// SelectTenant 选择租户
// @Tags 认证管理
// @Summary 选择租户
// @accept application/json
// @Produce application/json
// @Param req body dtoauth.SelectTenantReq true "选择租户请求"
// @Success 200 {object} gincontext.DtoRender{data=dtoauth.SelectTenantResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/auth/selectTenant [post]
func (ctr *authCtr) SelectTenant(ctx *gin.Context) {
	var req dtoauth.SelectTenantReq
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

// Logout 用户登出
// @Tags 认证管理
// @Summary 用户登出
// @accept application/json
// @Produce application/json
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "登出成功"}"
// @Router /v1/iam/auth/logout [post]
func (ctr *authCtr) Logout(ctx *gin.Context) {
	if err := ctr.authSvc.Logout(ctx); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "登出成功")
}
