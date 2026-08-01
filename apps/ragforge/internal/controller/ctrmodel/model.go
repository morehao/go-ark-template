package ctrmodel

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtomodel"
	"github.com/morehao/goark/ragforge/internal/service/svcmodel"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type ModelCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Test(ctx *gin.Context)
	GetProviders(ctx *gin.Context)
}

type modelCtr struct {
	modelSvc svcmodel.ModelSvc
}

var _ ModelCtr = (*modelCtr)(nil)

func NewModelCtr() ModelCtr {
	return &modelCtr{
		modelSvc: svcmodel.NewModelSvc(),
	}
}

func (ctr *modelCtr) Create(ctx *gin.Context) {
	var req dtomodel.ModelCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.modelSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *modelCtr) Delete(ctx *gin.Context) {
	var req dtomodel.ModelDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.modelSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

func (ctr *modelCtr) Update(ctx *gin.Context) {
	var req dtomodel.ModelUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.modelSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

func (ctr *modelCtr) Detail(ctx *gin.Context) {
	var req dtomodel.ModelDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.modelSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *modelCtr) PageList(ctx *gin.Context) {
	var req dtomodel.ModelPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.modelSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *modelCtr) Test(ctx *gin.Context) {
	var req dtomodel.ModelTestReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.modelSvc.Test(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *modelCtr) GetProviders(ctx *gin.Context) {
	res, err := ctr.modelSvc.GetProviders(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
