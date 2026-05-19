package chatpipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/morehao/golib/glog"
)

type PluginBuildContext struct{}

func NewPluginBuildContext(em *EventManager) *PluginBuildContext {
	p := &PluginBuildContext{}
	em.Register(p)
	return p
}

func (p *PluginBuildContext) ActivationEvents() []EventType {
	return []EventType{EventBuildContext}
}

func (p *PluginBuildContext) OnEvent(ctx context.Context, _ EventType, cm *ChatManage, next func() *PluginError) *PluginError {
	var sb strings.Builder

	if len(cm.History) > 0 {
		sb.WriteString("--- 对话历史 ---\n")
		for _, h := range cm.History {
			sb.WriteString(fmt.Sprintf("用户: %s\n助手: %s\n\n", h.Query, h.Answer))
		}
		sb.WriteString("--- 对话历史结束 ---\n\n")
	}

	if len(cm.SearchResults) > 0 {
		sb.WriteString("--- 相关上下文 ---\n")
		for i, r := range cm.SearchResults {
			sb.WriteString(fmt.Sprintf("[片段 %d]:\n%s\n\n", i+1, r.Content))
		}
		sb.WriteString("--- 相关上下文结束 ---\n\n")
	}

	sb.WriteString(fmt.Sprintf("用户问题: %s", cm.Query))
	cm.ContextContent = sb.String()

	glog.Infof(ctx, "[PluginBuildContext] context built, history:%d, searchResults:%d, len:%d",
		len(cm.History), len(cm.SearchResults), len(cm.ContextContent))
	return next()
}
