package ctrfaq

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtofaq"
	"github.com/morehao/goark/ragforge/internal/service/svcfaq"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type FAQCtr interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Search(ctx *gin.Context)
	Import(ctx *gin.Context)
}

type faqCtr struct {
	faqSvc svcfaq.FAQSvc
}

var _ FAQCtr = (*faqCtr)(nil)

func NewFAQCtr() FAQCtr {
	return &faqCtr{
		faqSvc: svcfaq.NewFAQSvc(),
	}
}

func (ctr *faqCtr) Create(ctx *gin.Context) {
	var req dtofaq.FAQCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.faqSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *faqCtr) Update(ctx *gin.Context) {
	var req dtofaq.FAQUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.faqSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

func (ctr *faqCtr) Delete(ctx *gin.Context) {
	var req dtofaq.FAQDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.faqSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

func (ctr *faqCtr) Detail(ctx *gin.Context) {
	var req dtofaq.FAQDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.faqSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *faqCtr) PageList(ctx *gin.Context) {
	var req dtofaq.FAQPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.faqSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *faqCtr) Search(ctx *gin.Context) {
	var req dtofaq.FAQSearchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.faqSvc.Search(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *faqCtr) Import(ctx *gin.Context) {
	gincontext.Success(ctx, "功能开发中")
}
