package ctrsession

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtosession"
	"github.com/morehao/goark/ragforge/internal/service/svcsession"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type SessionCtr interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	Detail(ctx *gin.Context)
	PageList(ctx *gin.Context)
	GenerateTitle(ctx *gin.Context)
	Stop(ctx *gin.Context)
}

type sessionCtr struct {
	sessionSvc svcsession.SessionSvc
}

var _ SessionCtr = (*sessionCtr)(nil)

func NewSessionCtr() SessionCtr {
	return &sessionCtr{
		sessionSvc: svcsession.NewSessionSvc(),
	}
}

// Create 创建会话
// @Tags 会话管理
// @Summary 创建会话
// @accept application/json
// @Produce application/json
// @Param req body dtosession.SessionCreateReq true "创建会话"
// @Success 200 {object} gincontext.DtoRender{data=dtosession.SessionCreateResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/session/create [post]
func (ctr *sessionCtr) Create(ctx *gin.Context) {
	var req dtosession.SessionCreateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.sessionSvc.Create(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *sessionCtr) Delete(ctx *gin.Context) {
	var req dtosession.SessionDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.sessionSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}

func (ctr *sessionCtr) Update(ctx *gin.Context) {
	var req dtosession.SessionUpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.sessionSvc.Update(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "修改成功")
	}
}

// Detail 会话详情
// @Tags 会话管理
// @Summary 会话详情
// @accept application/json
// @Produce application/json
// @Param req query dtosession.SessionDetailReq true "会话详情"
// @Success 200 {object} gincontext.DtoRender{data=dtosession.SessionDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/session/detail [get]
func (ctr *sessionCtr) Detail(ctx *gin.Context) {
	var req dtosession.SessionDetailReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.sessionSvc.Detail(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// GenerateTitle 生成会话标题
// @Tags 会话管理
// @Summary 生成会话标题
// @accept application/json
// @Produce application/json
// @Param req body dtosession.SessionGenerateTitleReq true "生成会话标题"
// @Success 200 {object} gincontext.DtoRender{data=dtosession.SessionGenerateTitleResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/session/generateTitle [post]
func (ctr *sessionCtr) GenerateTitle(ctx *gin.Context) {
	var req dtosession.SessionGenerateTitleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.sessionSvc.GenerateTitle(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Stop 停止会话
// @Tags 会话管理
// @Summary 停止会话
// @accept application/json
// @Produce application/json
// @Param req body dtosession.SessionStopReq true "停止会话"
// @Success 200 {object} gincontext.DtoRender{data=string} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/session/stop [post]
func (ctr *sessionCtr) Stop(ctx *gin.Context) {
	var req dtosession.SessionStopReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.sessionSvc.Stop(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "功能开发中")
	}
}

// PageList 会话列表
// @Tags 会话管理
// @Summary 会话列表分页
// @accept application/json
// @Produce application/json
// @Param req body dtosession.SessionPageListReq true "会话列表"
// @Success 200 {object} gincontext.DtoRender{data=dtosession.SessionPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/session/pageList [post]
func (ctr *sessionCtr) PageList(ctx *gin.Context) {
	var req dtosession.SessionPageListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.sessionSvc.PageList(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}
