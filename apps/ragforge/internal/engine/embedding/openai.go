package embedding

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"github.com/morehao/goark/ragforge/internal/engine"
)

type openAIEmbeddingProvider struct {
	client    *openai.Client
	modelName string
}

func NewOpenAIEmbeddingProvider(apiKey, apiBase, modelName string) engine.EmbeddingProvider {
	config := openai.DefaultConfig(apiKey)
	if apiBase != "" {
		config.BaseURL = apiBase
	}
	client := openai.NewClientWithConfig(config)
	return &openAIEmbeddingProvider{
		client:    client,
		modelName: modelName,
	}
}

func (p *openAIEmbeddingProvider) CreateEmbedding(ctx context.Context, req *engine.EmbeddingRequest) (*engine.EmbeddingResponse, error) {
	model := req.Model
	if model == "" {
		model = p.modelName
	}

	openaiReq := openai.EmbeddingRequestStrings{
		Model: openai.EmbeddingModel(model),
		Input: req.Inputs,
	}

	resp, err := p.client.CreateEmbeddings(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("openai create embedding: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for _, data := range resp.Data {
		embeddings[data.Index] = data.Embedding
	}

	return &engine.EmbeddingResponse{
		Embeddings: embeddings,
	}, nil
}
