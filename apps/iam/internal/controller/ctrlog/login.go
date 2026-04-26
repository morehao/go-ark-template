package ctrlog

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/service/svclog"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type loginLogCtr struct {
	loginLogSvc svclog.LoginLogSvc
}

func NewLoginLogCtr() *loginLogCtr {
	return &loginLogCtr{
		loginLogSvc: svclog.NewLoginLogSvc(),
	}
}

func (ctr *loginLogCtr) Create(ctx *gin.Context) {
	var req svclog.LoginLogCreateReq
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
	var req svclog.LoginLogPageListReq
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
