package ctrorg

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtoorg"
	"github.com/morehao/goark/apps/iam/internal/service/svcorg"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type DepartmentCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Tree(ctx *gin.Context)
}

type departmentCtr struct {
	departmentSvc svcorg.DepartmentSvc
}

var _ DepartmentCtr = (*departmentCtr)(nil)

func NewDepartmentCtr() DepartmentCtr {
	return &departmentCtr{
		departmentSvc: svcorg.NewDepartmentSvc(),
	}
}

// Create 创建部门管理
// @Tags 部门管理
// @Summary 创建部门管理
// @accept application/json
// @Produce application/json
// @Param req body dtoorg.DepartmentCreateReq true "创建部门管理"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.DepartmentCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/department/create [post]
func (ctr *departmentCtr) Create(ctx *gin.Context) {
	var req dtoorg.DepartmentCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.departmentSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Delete 删除部门管理
// @Tags 部门管理
// @Summary 删除部门管理
// @accept application/json
// @Produce application/json
// @Param req body dtoorg.DepartmentDeleteReq true "删除部门管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/iam/department/delete [post]
func (ctr *departmentCtr) Delete(ctx *gin.Context) {
	var req dtoorg.DepartmentDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.departmentSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

// Update 修改部门管理
// @Tags 部门管理
// @Summary 修改部门管理
// @accept application/json
// @Produce application/json
// @Param req body dtoorg.DepartmentUpdateReq true "修改部门管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/iam/department/update [post]
func (ctr *departmentCtr) Update(ctx *gin.Context) {
	var req dtoorg.DepartmentUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.departmentSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

// Detail 部门管理详情
// @Tags 部门管理
// @Summary 部门管理详情
// @accept application/json
// @Produce application/json
// @Param req query dtoorg.DepartmentDetailReq true "部门管理详情"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.DepartmentDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/department/detail [get]
func (ctr *departmentCtr) Detail(ctx *gin.Context) {
	var req dtoorg.DepartmentDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.departmentSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// PageList 部门管理列表
// @Tags 部门管理
// @Summary 部门管理列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtoorg.DepartmentPageListReq true "部门管理列表"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.DepartmentPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/department/pageList [post]
func (ctr *departmentCtr) PageList(ctx *gin.Context) {
	var req dtoorg.DepartmentPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.departmentSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Tree 获取部门树
// @Tags 部门管理
// @Summary 获取部门树
// @accept application/json
// @Produce application/json
// @Param req query dtoorg.DepartmentTreeReq true "部门树"
// @Success 200 {object} gincontext.DtoRender{data=dtoorg.DepartmentTreeResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/department/tree [get]
func (ctr *departmentCtr) Tree(ctx *gin.Context) {
	var req dtoorg.DepartmentTreeReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.departmentSvc.Tree(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
