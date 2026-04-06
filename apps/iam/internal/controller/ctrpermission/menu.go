package ctrpermission

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/iam/internal/dto/dtopermission"
	"github.com/morehao/goark/apps/iam/internal/service/svcpermission"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type MenuCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Tree(ctx *gin.Context)
}

type menuCtr struct {
	menuSvc svcpermission.MenuSvc
}

var _ MenuCtr = (*menuCtr)(nil)

func NewMenuCtr() MenuCtr {
	return &menuCtr{
		menuSvc: svcpermission.NewMenuSvc(),
	}
}

// Create 创建菜单管理
// @Tags 菜单管理
// @Summary 创建菜单管理
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.MenuCreateReq true "创建菜单管理"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.MenuCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/menu/create [post]
func (ctr *menuCtr) Create(ctx *gin.Context) {
	var req dtopermission.MenuCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.menuSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Delete 删除菜单管理
// @Tags 菜单管理
// @Summary 删除菜单管理
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.MenuDeleteReq true "删除菜单管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/iam/menu/delete [post]
func (ctr *menuCtr) Delete(ctx *gin.Context) {
	var req dtopermission.MenuDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	if err := ctr.menuSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

// Update 修改菜单管理
// @Tags 菜单管理
// @Summary 修改菜单管理
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.MenuUpdateReq true "修改菜单管理"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/iam/menu/update [post]
func (ctr *menuCtr) Update(ctx *gin.Context) {
	var req dtopermission.MenuUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.menuSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

// Detail 菜单管理详情
// @Tags 菜单管理
// @Summary 菜单管理详情
// @accept application/json
// @Produce application/json
// @Param req query dtopermission.MenuDetailReq true "菜单管理详情"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.MenuDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/menu/detail [get]
func (ctr *menuCtr) Detail(ctx *gin.Context) {
	var req dtopermission.MenuDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.menuSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// PageList 菜单管理列表
// @Tags 菜单管理
// @Summary 菜单管理列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtopermission.MenuPageListReq true "菜单管理列表"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.MenuPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/menu/pageList [post]
func (ctr *menuCtr) PageList(ctx *gin.Context) {
	var req dtopermission.MenuPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.menuSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Tree 获取菜单树
// @Tags 菜单管理
// @Summary 获取菜单树
// @accept application/json
// @Produce application/json
// @Param req query dtopermission.MenuTreeReq true "菜单树"
// @Success 200 {object} gincontext.DtoRender{data=dtopermission.MenuTreeResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/iam/menu/tree [get]
func (ctr *menuCtr) Tree(ctx *gin.Context) {
	var req dtopermission.MenuTreeReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.menuSvc.Tree(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
