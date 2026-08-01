package ctrkb

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtokb"
	"github.com/morehao/goark/ragforge/internal/service/svckb"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type KBCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Copy(ctx *gin.Context)
}

type kbCtr struct {
	kbSvc svckb.KBSvc
}

var _ KBCtr = (*kbCtr)(nil)

func NewKBCtr() KBCtr {
	return &kbCtr{
		kbSvc: svckb.NewKBSvc(),
	}
}

// Create 创建知识库
// @Tags 知识库管理
// @Summary 创建知识库
// @accept application/json
// @Produce application/json
// @Param req body dtokb.KBCreateReq true "创建知识库"
// @Success 200 {object} gincontext.DtoRender{data=dtokb.KBCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/kb/create [post]
func (ctr *kbCtr) Create(ctx *gin.Context) {
	var req dtokb.KBCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.kbSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Delete 删除知识库
// @Tags 知识库管理
// @Summary 删除知识库
// @accept application/json
// @Produce application/json
// @Param req body dtokb.KBDeleteReq true "删除知识库"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/ragforge/kb/delete [post]
func (ctr *kbCtr) Delete(ctx *gin.Context) {
	var req dtokb.KBDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.kbSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

// Update 修改知识库
// @Tags 知识库管理
// @Summary 修改知识库
// @accept application/json
// @Produce application/json
// @Param req body dtokb.KBUpdateReq true "修改知识库"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/ragforge/kb/update [post]
func (ctr *kbCtr) Update(ctx *gin.Context) {
	var req dtokb.KBUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.kbSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

// Detail 知识库详情
// @Tags 知识库管理
// @Summary 知识库详情
// @accept application/json
// @Produce application/json
// @Param req query dtokb.KBDetailReq true "知识库详情"
// @Success 200 {object} gincontext.DtoRender{data=dtokb.KBDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/kb/detail [get]
func (ctr *kbCtr) Detail(ctx *gin.Context) {
	var req dtokb.KBDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.kbSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Copy 复制知识库
// @Tags 知识库管理
// @Summary 复制知识库
// @accept application/json
// @Produce application/json
// @Param req body dtokb.KBCopyReq true "复制知识库"
// @Success 200 {object} gincontext.DtoRender{data=dtokb.KBCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/kb/copy [post]
func (ctr *kbCtr) Copy(ctx *gin.Context) {
	var req dtokb.KBCopyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.kbSvc.Copy(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// PageList 知识库列表
// @Tags 知识库管理
// @Summary 知识库列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtokb.KBPageListReq true "知识库列表"
// @Success 200 {object} gincontext.DtoRender{data=dtokb.KBPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/kb/pageList [post]
func (ctr *kbCtr) PageList(ctx *gin.Context) {
	var req dtokb.KBPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.kbSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
