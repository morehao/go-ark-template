package ctrauth

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoauth"
	"github.com/morehao/goark/apps/iam/internal/service/svcauth"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type LoginLogCtr interface {
	Create(ctx *gin.Context)
	PageList(ctx *gin.Context)
}

type loginLogCtr struct {
	loginLogSvc svcauth.LoginLogSvc
}

var _ LoginLogCtr = (*loginLogCtr)(nil)

func NewLoginLogCtr() LoginLogCtr {
	return &loginLogCtr{
		loginLogSvc: svcauth.NewLoginLogSvc(),
	}
}

func (ctr *loginLogCtr) Create(ctx *gin.Context) {
	var req dtoauth.LoginLogCreateReq
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

func (ctr *loginLogCtr) PageList(ctx *gin.Context) {
	var req dtoauth.LoginLogPageListReq
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

type OperationLogCtr interface {
	Create(ctx *gin.Context)
	PageList(ctx *gin.Context)
}

type operationLogCtr struct {
	operationLogSvc svcauth.OperationLogSvc
}

var _ OperationLogCtr = (*operationLogCtr)(nil)

func NewOperationLogCtr() OperationLogCtr {
	return &operationLogCtr{
		operationLogSvc: svcauth.NewOperationLogSvc(),
	}
}

func (ctr *operationLogCtr) Create(ctx *gin.Context) {
	var req dtoauth.OperationLogCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.operationLogSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *operationLogCtr) PageList(ctx *gin.Context) {
	var req dtoauth.OperationLogPageListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.operationLogSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}