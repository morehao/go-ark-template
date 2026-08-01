package chatpipeline

import (
	"context"
	"strings"

	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/golib/glog"
)

type PluginChatCompleteStream struct{}

func NewPluginChatCompleteStream(em *EventManager) *PluginChatCompleteStream {
	p := &PluginChatCompleteStream{}
	em.Register(p)
	return p
}

func (p *PluginChatCompleteStream) ActivationEvents() []EventType {
	return []EventType{EventChatCompleteStream}
}

func (p *PluginChatCompleteStream) OnEvent(ctx context.Context, _ EventType, cm *ChatManage, next func() *PluginError) *PluginError {
	if cm.LLMProvider == nil {
		glog.Errorf(ctx, "[PluginChatCompleteStream] LLM provider not available")
		if cm.OnDone != nil {
			cm.OnDone(buildFallbackAnswer(cm.Query, cm.SearchResults))
		}
		return next()
	}

	messages := []engine.ChatMessage{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: cm.ContextContent},
	}

	msgStream, err := cm.LLMProvider.ChatCompletionStream(ctx, &engine.ChatCompletionRequest{
		Messages:    messages,
		Temperature: 0.3,
		Stream:      true,
	})
	if err != nil {
		glog.Errorf(ctx, "[PluginChatCompleteStream] ChatCompletionStream fail, err:%v", err)
		if cm.OnDone != nil {
			cm.OnDone(buildFallbackAnswer(cm.Query, cm.SearchResults))
		}
		return next()
	}

	var fullContent strings.Builder
	for token := range msgStream {
		fullContent.WriteString(token)
		if cm.OnToken != nil {
			cm.OnToken(token)
		}
	}

	cm.Answer = fullContent.String()
	if cm.OnDone != nil {
		cm.OnDone(cm.Answer)
	}
	glog.Infof(ctx, "[PluginChatCompleteStream] answer generated, len:%d", len(cm.Answer))
	return next()
}
