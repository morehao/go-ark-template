package llm

import (
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"

	"github.com/morehao/goark/ragforge/internal/engine"
)

type openAIProvider struct {
	client *openai.Client
}

func NewOpenAIProvider(apiKey, apiBase string) engine.LLMProvider {
	config := openai.DefaultConfig(apiKey)
	if apiBase != "" {
		config.BaseURL = apiBase
	}
	client := openai.NewClientWithConfig(config)
	return &openAIProvider{
		client: client,
	}
}

func (p *openAIProvider) ChatCompletion(ctx context.Context, req *engine.ChatCompletionRequest) (*engine.ChatCompletionResponse, error) {
	openaiReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}

	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := p.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("openai chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("openai chat completion: no choices returned")
	}

	return &engine.ChatCompletionResponse{
		Content: resp.Choices[0].Message.Content,
	}, nil
}

func (p *openAIProvider) ChatCompletionStream(ctx context.Context, req *engine.ChatCompletionRequest) (<-chan string, error) {
	openaiReq := openai.ChatCompletionRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      true,
	}

	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("openai chat completion stream: %w", err)
	}

	ch := make(chan string)
	go func() {
		defer func() { _ = stream.Close() }()
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				return
			}
			if len(resp.Choices) > 0 {
				ch <- resp.Choices[0].Delta.Content
			}
		}
	}()

	return ch, nil
}
