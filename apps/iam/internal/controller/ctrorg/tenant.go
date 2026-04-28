package ctrorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/internal/service/svcorg"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type TenantCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
}

type tenantCtr struct {
	tenantSvc svcorg.TenantSvc
}

var _ TenantCtr = (*tenantCtr)(nil)

func NewTenantCtr() TenantCtr {
	return &tenantCtr{
		tenantSvc: svcorg.NewTenantSvc(),
	}
}

func (ctr *tenantCtr) Create(ctx *gin.Context) {
	var req dtoorg.TenantCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.tenantSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *tenantCtr) Delete(ctx *gin.Context) {
	var req dtoorg.TenantDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.tenantSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

func (ctr *tenantCtr) Update(ctx *gin.Context) {
	var req dtoorg.TenantUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.tenantSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

func (ctr *tenantCtr) Detail(ctx *gin.Context) {
	var req dtoorg.TenantDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.tenantSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *tenantCtr) PageList(ctx *gin.Context) {
	var req dtoorg.TenantPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.tenantSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
