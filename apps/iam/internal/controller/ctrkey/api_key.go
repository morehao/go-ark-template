package ctrkey

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoapikey"
	"github.com/morehao/goark/apps/iam/internal/service/svcapikey"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type apiKeyCtr struct {
	apiKeySvc svcapikey.ApiKeySvc
}

func NewApiKeyCtr() *apiKeyCtr {
	return &apiKeyCtr{
		apiKeySvc: svcapikey.NewApiKeySvc(),
	}
}

func (ctr *apiKeyCtr) Create(ctx *gin.Context) {
	var req dtoapikey.ApiKeyCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.apiKeySvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *apiKeyCtr) Delete(ctx *gin.Context) {
	var req dtoapikey.ApiKeyDeleteReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	err := ctr.apiKeySvc.Delete(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, nil)
}

func (ctr *apiKeyCtr) List(ctx *gin.Context) {
	var req dtoapikey.ApiKeyListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.apiKeySvc.List(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *apiKeyCtr) Disable(ctx *gin.Context) {
	var req dtoapikey.ApiKeyDisableReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	err := ctr.apiKeySvc.Disable(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, nil)
}

func (ctr *apiKeyCtr) Enable(ctx *gin.Context) {
	var req dtoapikey.ApiKeyEnableReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	err := ctr.apiKeySvc.Enable(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, nil)
}
