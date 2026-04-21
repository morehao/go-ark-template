package ctrorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/internal/service/svcorg"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type OrgCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetConfigsByDomain(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	ListConfig(ctx *gin.Context)
}

type orgCtr struct {
	orgSvc svcorg.OrgSvc
}

var _ OrgCtr = (*orgCtr)(nil)

func NewOrgCtr() OrgCtr {
	return &orgCtr{
		orgSvc: svcorg.NewOrgSvc(),
	}
}

func (ctr *orgCtr) Create(ctx *gin.Context) {
	var req dtoorg.OrgCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.orgSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *orgCtr) Delete(ctx *gin.Context) {
	var req dtoorg.OrgDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.orgSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

func (ctr *orgCtr) Update(ctx *gin.Context) {
	var req dtoorg.OrgUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.orgSvc.Update(ctx, &req); err != nil {
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
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.OrgConfigsResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/org/getConfigsByDomain [get]
func (ctr *orgCtr) GetConfigsByDomain(ctx *gin.Context) {
	var req dtoorg.OrgGetConfigsByDomainReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.orgSvc.GetConfigsByDomain(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *orgCtr) Detail(ctx *gin.Context) {
	var req dtoorg.OrgDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.orgSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *orgCtr) PageList(ctx *gin.Context) {
	var req dtoorg.OrgPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.orgSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *orgCtr) ListConfig(ctx *gin.Context) {
	res, err := ctr.orgSvc.ListConfig(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}