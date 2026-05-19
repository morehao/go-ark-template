package chatpipeline

import (
	"context"

	"github.com/morehao/goark/ragforge/internal/engine"
)

const (
	SystemPrompt = "你是一个专业的知识库助手。请根据提供的上下文内容回答用户问题。如果你不确定答案，请如实说不知道。"

	FallbackMsgWithResults = "基于知识库内容，找到以下相关信息:\n\n%s\n\n(提示: 请配置 LLM 模型以获得更好的回答质量)"
	FallbackMsgNoResults   = "暂未找到相关知识库内容。请确保知识库已添加文档并完成向量化。"
)

type EventType string

const (
	EventLoadHistory       EventType = "load_history"
	EventSearch            EventType = "search"
	EventBuildContext      EventType = "build_context"
	EventChatComplete      EventType = "chat_complete"
	EventChatCompleteStream EventType = "chat_complete_stream"
)

type PluginError struct {
	Err         error
	Description string
	ErrorType   string
}

type ChatManage struct {
	SessionID         uint
	KbID              uint
	TenantID          uint
	UserID            uint
	Query             string
	History           []History
	SearchResults     []SearchResult
	ContextContent    string
	Answer            string
	LLMProvider       engine.LLMProvider
	EmbeddingProvider engine.EmbeddingProvider
	OnToken           func(token string)
	OnDone            func(answer string)
}

type History struct {
	Query  string
	Answer string
}

type SearchResult struct {
	Content string
	Score   float64
	ChunkID uint
}

type Plugin interface {
	OnEvent(ctx context.Context, eventType EventType, chatManage *ChatManage, next func() *PluginError) *PluginError
	ActivationEvents() []EventType
}
