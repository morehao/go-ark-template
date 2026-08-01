package chatpipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/morehao/goark/ragforge/internal/engine"
	"github.com/morehao/golib/glog"
)

type PluginChatComplete struct{}

func NewPluginChatComplete(em *EventManager) *PluginChatComplete {
	p := &PluginChatComplete{}
	em.Register(p)
	return p
}

func (p *PluginChatComplete) ActivationEvents() []EventType {
	return []EventType{EventChatComplete}
}

func (p *PluginChatComplete) OnEvent(ctx context.Context, _ EventType, cm *ChatManage, next func() *PluginError) *PluginError {
	if cm.LLMProvider == nil {
		glog.Errorf(ctx, "[PluginChatComplete] LLM provider not available")
		cm.Answer = buildFallbackAnswer(cm.Query, cm.SearchResults)
		return next()
	}

	messages := []engine.ChatMessage{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: cm.ContextContent},
	}

	resp, err := cm.LLMProvider.ChatCompletion(ctx, &engine.ChatCompletionRequest{
		Messages:    messages,
		Temperature: 0.3,
	})
	if err != nil {
		glog.Errorf(ctx, "[PluginChatComplete] ChatCompletion fail, err:%v", err)
		cm.Answer = buildFallbackAnswer(cm.Query, cm.SearchResults)
		return next()
	}

	cm.Answer = resp.Content
	glog.Infof(ctx, "[PluginChatComplete] answer generated, len:%d", len(cm.Answer))
	return next()
}

func buildFallbackAnswer(_ string, results []SearchResult) string {
	if len(results) > 0 {
		var sb strings.Builder
		sb.WriteString("基于知识库内容，找到以下相关信息:\n\n")
		for i, r := range results {
			fmt.Fprintf(&sb, "%d. %s\n", i+1, r.Content)
		}
		sb.WriteString("\n(提示: 请配置 LLM 模型以获得更好的回答质量)")
		return sb.String()
	}
	return FallbackMsgNoResults
}
