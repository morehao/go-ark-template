package ctrknowledge

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/ragflow/internal/dto/dtoknowledge"
	"github.com/morehao/goark/apps/ragflow/internal/service/svcknowledge"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type KnowledgeBaseCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	List(ctx *gin.Context)
}

type knowledgeBaseCtr struct {
	knowledgeBaseSvc svcknowledge.KnowledgeBaseSvc
}

var _ KnowledgeBaseCtr = (*knowledgeBaseCtr)(nil)

func NewKnowledgeBaseCtr() KnowledgeBaseCtr {
	return &knowledgeBaseCtr{
		knowledgeBaseSvc: svcknowledge.NewKnowledgeBaseSvc(),
	}
}

func (ctr *knowledgeBaseCtr) Create(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeBaseCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeBaseSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *knowledgeBaseCtr) Delete(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeBaseDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.knowledgeBaseSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

func (ctr *knowledgeBaseCtr) Update(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeBaseUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.knowledgeBaseSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

func (ctr *knowledgeBaseCtr) Detail(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeBaseDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeBaseSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *knowledgeBaseCtr) List(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeBaseListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeBaseSvc.List(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}