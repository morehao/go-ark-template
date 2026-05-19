package ctrqa

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/ragforge/internal/dto/dtoqa"
	"github.com/morehao/goark/ragforge/internal/service/svcqa"
	"github.com/morehao/golib/biz/gcontext/gincontext"
)

type QACtr interface {
	KnowledgeChat(ctx *gin.Context)
	KnowledgeChatStream(ctx *gin.Context)
}

type qaCtr struct {
	qaSvc svcqa.QASvc
}

var _ QACtr = (*qaCtr)(nil)

func NewQACtr() QACtr {
	return &qaCtr{
		qaSvc: svcqa.NewQASvc(),
	}
}

// KnowledgeChat 知识库问答
// @Tags 问答管理
// @Summary 知识库问答
// @accept application/json
// @Produce application/json
// @Param req body dtoqa.KnowledgeChatReq true "知识库问答"
// @Success 200 {object} gincontext.DtoRender{data=dtoqa.KnowledgeChatResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /v1/ragforge/qa/knowledgeChat [post]
func (ctr *qaCtr) KnowledgeChat(ctx *gin.Context) {
	var req dtoqa.KnowledgeChatReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}
	res, err := ctr.qaSvc.KnowledgeChat(ctx, &req)
	if err != nil {
		gincontext.Fail(ctx, err)
		return
	} else {
		gincontext.Success(ctx, res)
	}
}

type sseEvent struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	SessionID uint `json:"sessionId,omitempty"`
}

// KnowledgeChatStream 知识库问答流式
// @Tags 问答管理
// @Summary 知识库问答流式
// @accept application/json
// @Produce text/event-stream
// @Param req body dtoqa.KnowledgeChatReq true "知识库问答流式"
// @Success 200 {object} string "流式响应"
// @Router /v1/ragforge/qa/knowledgeChatStream [post]
func (ctr *qaCtr) KnowledgeChatStream(ctx *gin.Context) {
	var req dtoqa.KnowledgeChatReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gincontext.Fail(ctx, err)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.WriteHeader(200)

	flusher, ok := ctx.Writer.(gin.ResponseWriter)
	if !ok {
		return
	}

	writeEvent := func(event sseEvent) {
		data, _ := json.Marshal(event)
		fmt.Fprintf(ctx.Writer, "data: %s\n\n", data)
		flusher.Flush()
	}

	err := ctr.qaSvc.KnowledgeChatStream(ctx, &req,
		func(token string) {
			writeEvent(sseEvent{Type: "token", Content: token})
		},
		func(answer string) {
			writeEvent(sseEvent{Type: "done", Content: answer})
		},
	)
	if err != nil {
		writeEvent(sseEvent{Type: "error", Content: err.Error()})
	}
	
		ctx.Abort()
}
