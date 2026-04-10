package ctrauth

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/apps/iam/internal/service/svcauth"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type AuthCtr interface {
	LoginByPassword(ctx *gin.Context)
	SelectTenant(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
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

// LoginByPassword 密码登录
// @Tags 认证管理
// @Summary 密码登录
// @accept application/json
// @Produce application/json
// @Param req body dtoauth.LoginByPasswordReq true "密码登录请求"
// @Success 200 {object} gincontext.DtoRender{data=dtoauth.LoginByPasswordResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/auth/loginByPassword [post]
func (ctr *authCtr) LoginByPassword(ctx *gin.Context) {
	var req dtoauth.LoginByPasswordReq
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
// @Param req body dtoauth.LogoutReq true "登出请求"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "登出成功"}"
// @Router /v1/iam/auth/logout [post]
func (ctr *authCtr) Logout(ctx *gin.Context) {
	var req dtoauth.LogoutReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		req.RefreshToken = ""
	}
	if err := ctr.authSvc.Logout(ctx, req.RefreshToken); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "登出成功")
}

// RefreshToken 刷新令牌
// @Tags 认证管理
// @Summary 刷新令牌
// @accept application/json
// @Produce application/json
// @Param req body dtoauth.RefreshTokenReq true "刷新令牌请求"
// @Success 200 {object} gincontext.DtoRender{data=dtoauth.RefreshTokenResp} "{"code": 0,"data": {"token":"...","refreshToken":"..."},"msg": "success"}"
// @Router /v1/iam/auth/refreshToken [post]
func (ctr *authCtr) RefreshToken(ctx *gin.Context) {
	var req dtoauth.RefreshTokenReq
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
