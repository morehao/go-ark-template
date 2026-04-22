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
	GetOrgConfig(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	ListConfigDefinitions(ctx *gin.Context)
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

// GetOrgConfig 根据域名获取组织配置
// @Tags 组织管理
// @Summary 根据域名获取组织配置
// @accept application/json
// @Produce application/json
// @Param domain query dtoorg.GetOrganizationConfigsReq true "获取组织配置"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.GetOrgConfigResp} "{"code": 0, "requestID": "xxx", "data": "ok", "msg": "success"}"
// @Router /v1/iam/organization/getOrgConfig [get]
func (ctr *organizationCtr) GetOrgConfig(ctx *gin.Context) {
	var req dtoorg.GetOrganizationConfigsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.orgSvc.GetOrgConfig(ctx, &req)
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

// ListConfigDefinitions 获取配置项定义列表
// @Tags 组织管理
// @Summary 获取配置项定义列表
// @accept application/json
// @Produce application/json
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.ListConfigDefinitionsResp} "{"code": 0, "requestID": "xxx", "data": "ok", "msg": "success"}"
// @Router /v1/iam/organization/listConfigDefinitions [get]
func (ctr *organizationCtr) ListConfigDefinitions(ctx *gin.Context) {
	res, err := ctr.orgSvc.ListConfigDefinitions(ctx)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
