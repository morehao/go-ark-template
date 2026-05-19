package ctrknowledge

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtoknowledge"
	"github.com/morehao/goark/ragforge/internal/service/svcknowledge"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type KnowledgeCtr interface {
	CreateFromFile(ctx *gin.Context)
	CreateFromURL(ctx *gin.Context)
	CreateManual(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	Reparse(ctx *gin.Context)
	Download(ctx *gin.Context)
}

type knowledgeCtr struct {
	knowledgeSvc svcknowledge.KnowledgeSvc
}

var _ KnowledgeCtr = (*knowledgeCtr)(nil)

func NewKnowledgeCtr() KnowledgeCtr {
	return &knowledgeCtr{
		knowledgeSvc: svcknowledge.NewKnowledgeSvc(),
	}
}

// CreateFromFile 从文件创建知识
// @Tags 知识管理
// @Summary 从文件创建知识
// @accept multipart/form-data
// @Produce application/json
// @Param kbId formData uint true "知识库ID"
// @Param title formData string true "标题"
// @Param file formData file true "文件"
// @Success 200 {object} gincontext.DtoRender{data=dtoknowledge.KnowledgeCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/knowledge/createFile [post]
func (ctr *knowledgeCtr) CreateFromFile(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeCreateFileReq
	if err := ctx.ShouldBind(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeSvc.CreateFromFile(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// CreateFromURL 从URL创建知识
// @Tags 知识管理
// @Summary 从URL创建知识
// @accept application/json
// @Produce application/json
// @Param req body dtoknowledge.KnowledgeCreateURLReq true "从URL创建知识"
// @Success 200 {object} gincontext.DtoRender{data=dtoknowledge.KnowledgeCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/knowledge/createUrl [post]
func (ctr *knowledgeCtr) CreateFromURL(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeCreateURLReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeSvc.CreateFromURL(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// CreateManual 手动创建知识
// @Tags 知识管理
// @Summary 手动创建知识
// @accept application/json
// @Produce application/json
// @Param req body dtoknowledge.KnowledgeCreateManualReq true "手动创建知识"
// @Success 200 {object} gincontext.DtoRender{data=dtoknowledge.KnowledgeCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/knowledge/createManual [post]
func (ctr *knowledgeCtr) CreateManual(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeCreateManualReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeSvc.CreateManual(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// Update 修改知识
// @Tags 知识管理
// @Summary 修改知识
// @accept application/json
// @Produce application/json
// @Param req body dtoknowledge.KnowledgeUpdateReq true "修改知识"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /v1/ragforge/knowledge/update [post]
func (ctr *knowledgeCtr) Update(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.knowledgeSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "修改成功")
}

// Delete 删除知识
// @Tags 知识管理
// @Summary 删除知识
// @accept application/json
// @Produce application/json
// @Param req body dtoknowledge.KnowledgeDeleteReq true "删除知识"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /v1/ragforge/knowledge/delete [post]
func (ctr *knowledgeCtr) Delete(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.knowledgeSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "删除成功")
}

// Detail 知识详情
// @Tags 知识管理
// @Summary 知识详情
// @accept application/json
// @Produce application/json
// @Param req query dtoknowledge.KnowledgeDetailReq true "知识详情"
// @Success 200 {object} gincontext.DtoRender{data=dtoknowledge.KnowledgeDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/knowledge/detail [get]
func (ctr *knowledgeCtr) Detail(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

// PageList 知识列表
// @Tags 知识管理
// @Summary 知识列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtoknowledge.KnowledgePageListReq true "知识列表"
// @Success 200 {object} gincontext.DtoRender{data=dtoknowledge.KnowledgePageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/knowledge/pageList [post]
func (ctr *knowledgeCtr) PageList(ctx *gin.Context) {
	var req dtoknowledge.KnowledgePageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.knowledgeSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, res)
}

func (ctr *knowledgeCtr) Download(ctx *gin.Context) {
	gincontext.Success(ctx, "功能开发中")
}

// Reparse 重试解析知识
// @Tags 知识管理
// @Summary 重试解析知识
// @accept application/json
// @Produce application/json
// @Param req body dtoknowledge.KnowledgeReparseReq true "重试解析知识"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "重试解析成功"}"
// @Router /v1/ragforge/knowledge/reparse [post]
func (ctr *knowledgeCtr) Reparse(ctx *gin.Context) {
	var req dtoknowledge.KnowledgeReparseReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.knowledgeSvc.Reparse(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	gincontext.Success(ctx, "重试解析成功")
}
