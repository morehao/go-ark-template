package llm

import (
	"github.com/morehao/goark/ragforge/internal/engine"
)

type LLMType string

const (
	LLMTypeOpenAI LLMType = "openai"
)

func NewLLMProvider(llmType LLMType, apiKey, apiBase string) (engine.LLMProvider, error) {
	switch llmType {
	case LLMTypeOpenAI:
		return NewOpenAIProvider(apiKey, apiBase), nil
	}
	return nil, nil
}
