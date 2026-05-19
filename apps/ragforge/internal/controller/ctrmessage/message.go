package ctrmessage

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtomessage"
	"github.com/morehao/goark/ragforge/internal/service/svcmessage"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type MessageCtr interface {
	List(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Search(ctx *gin.Context)
}

type messageCtr struct {
	messageSvc svcmessage.MessageSvc
}

var _ MessageCtr = (*messageCtr)(nil)

func NewMessageCtr() MessageCtr {
	return &messageCtr{
		messageSvc: svcmessage.NewMessageSvc(),
	}
}

func (ctr *messageCtr) List(ctx *gin.Context) {
	var req dtomessage.MessageListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.messageSvc.List(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

// Search 搜索消息
// @Tags 消息管理
// @Summary 搜索消息
// @accept application/json
// @Produce application/json
// @Param req body dtomessage.MessageSearchReq true "搜索消息"
// @Success 200 {object} gincontext.DtoRender{data=dtomessage.MessageSearchResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/message/search [post]
func (ctr *messageCtr) Search(ctx *gin.Context) {
	var req dtomessage.MessageSearchReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.messageSvc.Search(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

func (ctr *messageCtr) Delete(ctx *gin.Context) {
	var req dtomessage.MessageDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	if err := ctr.messageSvc.Delete(ctx, &req); err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, "删除成功")
	}
}
