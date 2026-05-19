package ctrchunk

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtochunk"
	"github.com/morehao/goark/ragforge/internal/service/svcchunk"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type ChunkCtr interface {
	PageList(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Search(ctx *gin.Context)
	Detail(ctx *gin.Context)
}

type chunkCtr struct {
	chunkSvc svcchunk.ChunkSvc
}

var _ ChunkCtr = (*chunkCtr)(nil)

func NewChunkCtr() ChunkCtr {
	return &chunkCtr{
		chunkSvc: svcchunk.NewChunkSvc(),
	}
}

// PageList 知识块列表
// @Tags 知识块管理
// @Summary 知识块列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtochunk.ChunkPageListReq true "知识块列表"
// @Success 200 {object} gincontext.DtoRender{data=dtochunk.ChunkPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/chunk/pageList [post]
func (ctr *chunkCtr) PageList(ctx *gin.Context) {
	var req dtochunk.ChunkPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.chunkSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Update 修改知识块
// @Tags 知识块管理
// @Summary 修改知识块
// @accept application/json
// @Produce application/json
// @Param req body dtochunk.ChunkUpdateReq true "修改知识块"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/ragforge/chunk/update [post]
func (ctr *chunkCtr) Update(ctx *gin.Context) {
	var req dtochunk.ChunkUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.chunkSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

// Delete 删除知识块
// @Tags 知识块管理
// @Summary 删除知识块
// @accept application/json
// @Produce application/json
// @Param req body dtochunk.ChunkDeleteReq true "删除知识块"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/ragforge/chunk/delete [post]
func (ctr *chunkCtr) Delete(ctx *gin.Context) {
	var req dtochunk.ChunkDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.chunkSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

// Detail 知识块详情
// @Tags 知识块管理
// @Summary 知识块详情
// @accept application/json
// @Produce application/json
// @Param req query dtochunk.ChunkDetailReq true "知识块详情"
// @Success 200 {object} gincontext.DtoRender{data=dtochunk.ChunkDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/chunk/detail [get]
func (ctr *chunkCtr) Detail(ctx *gin.Context) {
	var req dtochunk.ChunkDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.chunkSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Search 搜索知识块
// @Tags 知识块管理
// @Summary 搜索知识块
// @accept application/json
// @Produce application/json
// @Param req body dtochunk.ChunkSearchReq true "搜索知识块"
// @Success 200 {object} gincontext.DtoRender{data=dtochunk.ChunkSearchResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/chunk/search [post]
func (ctr *chunkCtr) Search(ctx *gin.Context) {
	var req dtochunk.ChunkSearchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.chunkSvc.Search(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}
