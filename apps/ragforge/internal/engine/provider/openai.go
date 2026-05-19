package provider

import (
	"fmt"
	"strings"
)

const (
	OpenAIBaseURL = "https://api.openai.com/v1"
)

type OpenAIProvider struct{}

func init() {
	Register(&OpenAIProvider{})
}

func (p *OpenAIProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderOpenAI,
		DisplayName: "OpenAI",
		Description: "gpt-5.2, gpt-5-mini, etc.",
		DefaultURLs: map[ModelType]string{
			ModelTypeChat:      OpenAIBaseURL,
			ModelTypeEmbedding: OpenAIBaseURL,
			ModelTypeRerank:    OpenAIBaseURL,
		},
		ModelTypes: []ModelType{
			ModelTypeChat,
			ModelTypeEmbedding,
			ModelTypeRerank,
		},
		RequiresAuth: true,
	}
}

func (p *OpenAIProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for OpenAI provider")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}

func IsOpenAIReasoningOrGPT5Model(modelName string) bool {
	name := strings.ToLower(strings.TrimSpace(modelName))
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, "gpt-5") {
		return true
	}
	for _, prefix := range []string{"o1", "o3", "o4"} {
		if name == prefix || strings.HasPrefix(name, prefix+"-") {
			return true
		}
	}
	return false
}
