package ctrvectorstore

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtovectorstore"
	"github.com/morehao/goark/ragforge/internal/service/svcvectorstore"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type VectorStoreCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Test(ctx *gin.Context)
	GetTypes(ctx *gin.Context)
}

type vectorStoreCtr struct {
	vectorStoreSvc svcvectorstore.VectorStoreSvc
}

var _ VectorStoreCtr = (*vectorStoreCtr)(nil)

func NewVectorStoreCtr() VectorStoreCtr {
	return &vectorStoreCtr{
		vectorStoreSvc: svcvectorstore.NewVectorStoreSvc(),
	}
}

func (ctr *vectorStoreCtr) Create(ctx *gin.Context) {
	var req dtovectorstore.VectorStoreCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.vectorStoreSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *vectorStoreCtr) Delete(ctx *gin.Context) {
	var req dtovectorstore.VectorStoreDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.vectorStoreSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

func (ctr *vectorStoreCtr) Update(ctx *gin.Context) {
	var req dtovectorstore.VectorStoreUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.vectorStoreSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

func (ctr *vectorStoreCtr) Detail(ctx *gin.Context) {
	var req dtovectorstore.VectorStoreDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.vectorStoreSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *vectorStoreCtr) PageList(ctx *gin.Context) {
	var req dtovectorstore.VectorStorePageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.vectorStoreSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *vectorStoreCtr) Test(ctx *gin.Context) {
	var req dtovectorstore.VectorStoreTestReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.vectorStoreSvc.Test(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *vectorStoreCtr) GetTypes(ctx *gin.Context) {
	res, err := ctr.vectorStoreSvc.GetTypes(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
