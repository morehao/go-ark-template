package svcqa

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/pkg/code"
	"github.com/morehao/goark/ragforge/dao"
	"github.com/morehao/goark/ragforge/internal/chatpipeline"
	"github.com/morehao/goark/ragforge/internal/dto/dtoqa"
	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/goark/ragforge/model"
	"github.com/morehao/golib/biz/gcontext/gincontext"
	"github.com/morehao/golib/glog"
)

type QASvc interface {
	KnowledgeChat(ctx *gin.Context, req *dtoqa.KnowledgeChatReq) (*dtoqa.KnowledgeChatResp, error)
	KnowledgeChatStream(ctx *gin.Context, req *dtoqa.KnowledgeChatReq, onToken func(string), onDone func(string)) error
}

type qaSvc struct {
}

var _ QASvc = (*qaSvc)(nil)

func NewQASvc() QASvc {
	return &qaSvc{}
}

func (svc *qaSvc) KnowledgeChat(ctx *gin.Context, req *dtoqa.KnowledgeChatReq) (*dtoqa.KnowledgeChatResp, error) {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	sessionID, err := svc.getOrCreateSession(ctx, req, userID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := svc.saveMessage(ctx, sessionID, tenantID, model.MessageRoleUser, req.Content); err != nil {
		return nil, err
	}

	engineFactory := engine.GetGlobalFactory()
	cm := &chatpipeline.ChatManage{
		SessionID:         sessionID,
		KbID:              req.KbID,
		TenantID:          tenantID,
		UserID:            userID,
		Query:             req.Content,
		LLMProvider:       engineFactory.GetLLM(),
		EmbeddingProvider: engineFactory.GetEmbedding(),
	}

	em := chatpipeline.NewEventManager()
	chatpipeline.NewPluginLoadHistory(em)
	chatpipeline.NewPluginSearch(em)
	chatpipeline.NewPluginBuildContext(em)
	chatpipeline.NewPluginChatComplete(em)

	for _, et := range []chatpipeline.EventType{
		chatpipeline.EventLoadHistory,
		chatpipeline.EventSearch,
		chatpipeline.EventBuildContext,
		chatpipeline.EventChatComplete,
	} {
		if err := em.Trigger(ctx, et, cm); err != nil {
			glog.Errorf(ctx, "[svcqa.KnowledgeChat] pipeline event %s fail, err:%v", et, err)
			break
		}
	}

	if err := svc.saveMessage(ctx, sessionID, tenantID, model.MessageRoleAssistant, cm.Answer); err != nil {
		return nil, err
	}

	return &dtoqa.KnowledgeChatResp{
		SessionID: sessionID,
		Content:   cm.Answer,
	}, nil
}

func (svc *qaSvc) KnowledgeChatStream(ctx *gin.Context, req *dtoqa.KnowledgeChatReq, onToken func(string), onDone func(string)) error {
	userID := gincontext.GetUserID(ctx)
	tenantID := gincontext.GetTenantID(ctx)

	sessionID, err := svc.getOrCreateSession(ctx, req, userID, tenantID)
	if err != nil {
		return err
	}

	if err := svc.saveMessage(ctx, sessionID, tenantID, model.MessageRoleUser, req.Content); err != nil {
		return err
	}

	engineFactory := engine.GetGlobalFactory()
	cm := &chatpipeline.ChatManage{
		SessionID:         sessionID,
		KbID:              req.KbID,
		TenantID:          tenantID,
		UserID:            userID,
		Query:             req.Content,
		LLMProvider:       engineFactory.GetLLM(),
		EmbeddingProvider: engineFactory.GetEmbedding(),
		OnToken:           onToken,
		OnDone: func(answer string) {
			if saveErr := svc.saveMessage(ctx, sessionID, tenantID, model.MessageRoleAssistant, answer); saveErr != nil {
				glog.Errorf(ctx, "[svcqa.KnowledgeChatStream] save assistant message fail, err:%v", saveErr)
			}
			onDone(answer)
		},
	}

	em := chatpipeline.NewEventManager()
	chatpipeline.NewPluginLoadHistory(em)
	chatpipeline.NewPluginSearch(em)
	chatpipeline.NewPluginBuildContext(em)
	chatpipeline.NewPluginChatCompleteStream(em)

	for _, et := range []chatpipeline.EventType{
		chatpipeline.EventLoadHistory,
		chatpipeline.EventSearch,
		chatpipeline.EventBuildContext,
		chatpipeline.EventChatCompleteStream,
	} {
		if err := em.Trigger(ctx, et, cm); err != nil {
			glog.Errorf(ctx, "[svcqa.KnowledgeChatStream] pipeline event %s fail, err:%v", et, err)
			break
		}
	}

	return nil
}

func (svc *qaSvc) getOrCreateSession(ctx *gin.Context, req *dtoqa.KnowledgeChatReq, userID, tenantID uint) (uint, error) {
	if req.SessionID > 0 {
		return req.SessionID, nil
	}
	sessionEntity := &model.SessionEntity{
		TenantID: tenantID,
		UserID:   userID,
		KbID:     req.KbID,
	}
	if err := dao.NewSessionDao().Insert(ctx, sessionEntity); err != nil {
		glog.Errorf(ctx, "[svcqa.getOrCreateSession] create session fail, err:%v", err)
		return 0, code.GetError(code.QAChatError)
	}
	title := req.Content
	if len([]rune(title)) > 100 {
		title = string([]rune(title)[:100])
	}
	_ = dao.NewSessionDao().UpdateByID(ctx, sessionEntity.ID, &model.SessionEntity{Title: title})
	return sessionEntity.ID, nil
}

func (svc *qaSvc) saveMessage(ctx *gin.Context, sessionID, tenantID uint, role model.MessageRole, content string) error {
	msg := &model.MessageEntity{
		SessionID: sessionID,
		TenantID:  tenantID,
		Role:      role,
		Content:   content,
	}
	if err := dao.NewMessageDao().Insert(ctx, msg); err != nil {
		glog.Errorf(ctx, "[svcqa.saveMessage] fail, err:%v", err)
		return code.GetError(code.QAChatError)
	}
	return nil
}
