package ctrlog

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtolog"
	"github.com/morehao/goark/apps/iam/internal/service/svclog"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type operationLogCtr struct {
	operationLogSvc svclog.OperationLogSvc
}

func NewOperationLogCtr() *operationLogCtr {
	return &operationLogCtr{
		operationLogSvc: svclog.NewOperationLogSvc(),
	}
}

func (ctr *operationLogCtr) Create(ctx *gin.Context) {
	var req dtolog.OperationLogCreateReq
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
	var req dtolog.OperationLogPageListReq
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
