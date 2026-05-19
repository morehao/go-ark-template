package ctrtag

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtotag"
	"github.com/morehao/goark/ragforge/internal/service/svctag"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type TagCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	List(ctx *gin.Context)
}

type tagCtr struct {
	tagSvc svctag.TagSvc
}

var _ TagCtr = (*tagCtr)(nil)

func NewTagCtr() TagCtr {
	return &tagCtr{
		tagSvc: svctag.NewTagSvc(),
	}
}

func (ctr *tagCtr) Create(ctx *gin.Context) {
	var req dtotag.TagCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.tagSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *tagCtr) Delete(ctx *gin.Context) {
	var req dtotag.TagDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.tagSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

func (ctr *tagCtr) Update(ctx *gin.Context) {
	var req dtotag.TagUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.tagSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

func (ctr *tagCtr) List(ctx *gin.Context) {
	var req dtotag.TagListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.tagSvc.List(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
