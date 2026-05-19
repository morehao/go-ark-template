package embedding

import (
	"github.com/morehao/goark/ragforge/internal/engine"
)

type EmbeddingType string

const (
	EmbeddingTypeOpenAI EmbeddingType = "openai"
)

func NewEmbeddingProvider(embType EmbeddingType, apiKey, apiBase, modelName string) (engine.EmbeddingProvider, error) {
	switch embType {
	case EmbeddingTypeOpenAI:
		return NewOpenAIEmbeddingProvider(apiKey, apiBase, modelName), nil
	}
	return nil, nil
}
