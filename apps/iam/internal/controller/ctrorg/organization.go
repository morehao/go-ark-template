package ctrorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/internal/service/svcorg"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type OrganizationCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetConfigsByDomain(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	ListConfig(ctx *gin.Context)
}

type organizationCtr struct {
	organizationSvc svcorg.OrganizationSvc
}

var _ OrganizationCtr = (*organizationCtr)(nil)

func NewOrganizationCtr() OrganizationCtr {
	return &organizationCtr{
		organizationSvc: svcorg.NewOrganizationSvc(),
	}
}

func (ctr *organizationCtr) Create(ctx *gin.Context) {
	var req dtoorg.OrganizationCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.organizationSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *organizationCtr) Delete(ctx *gin.Context) {
	var req dtoorg.OrganizationDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.organizationSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

func (ctr *organizationCtr) Update(ctx *gin.Context) {
	var req dtoorg.OrganizationUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.organizationSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

// GetConfigsByDomain 根据域名获取组织配置
// @Tags 组织管理
// @Summary 根据域名获取组织配置
// @accept application/json
// @Produce application/json
// @Param domain query string false "组织域名(可选)"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.OrganizationConfigsResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/organization/getConfigsByDomain [get]
func (ctr *organizationCtr) GetConfigsByDomain(ctx *gin.Context) {
	var req dtoorg.OrganizationGetConfigsByDomainReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.organizationSvc.GetConfigsByDomain(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *organizationCtr) Detail(ctx *gin.Context) {
	var req dtoorg.OrganizationDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.organizationSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *organizationCtr) PageList(ctx *gin.Context) {
	var req dtoorg.OrganizationPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.organizationSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *organizationCtr) ListConfig(ctx *gin.Context) {
	res, err := ctr.organizationSvc.ListConfig(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
