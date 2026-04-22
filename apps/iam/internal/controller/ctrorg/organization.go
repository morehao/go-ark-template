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
	GetConfigs(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	ListConfig(ctx *gin.Context)
}

type organizationCtr struct {
	orgSvc svcorg.OrganizationSvc
}

var _ OrganizationCtr = (*organizationCtr)(nil)

func NewOrganizationCtr() OrganizationCtr {
	return &organizationCtr{
		orgSvc: svcorg.NewOrganizationSvc(),
	}
}

func (ctr *organizationCtr) Create(ctx *gin.Context) {
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

func (ctr *organizationCtr) Delete(ctx *gin.Context) {
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

func (ctr *organizationCtr) Update(ctx *gin.Context) {
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

// GetConfigs 根据域名获取组织配置
// @Tags 组织管理
// @Summary 根据域名获取组织配置
// @accept application/json
// @Produce application/json
// @Param domain query dtoorg.GetOrganizationConfigsReq true "获取组织配置"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.GetOrganizationConfigsResp} "{"code": 0, "requestID": "xxx", "data": "ok", "msg": "success"}"
// @Router /v1/iam/organization/getConfigs [get]
func (ctr *organizationCtr) GetConfigs(ctx *gin.Context) {
	var req dtoorg.GetOrganizationConfigsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.orgSvc.GetConfigs(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *organizationCtr) Detail(ctx *gin.Context) {
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

func (ctr *organizationCtr) PageList(ctx *gin.Context) {
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

func (ctr *organizationCtr) ListConfig(ctx *gin.Context) {
	res, err := ctr.orgSvc.ListConfig(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
