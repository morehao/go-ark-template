package provider

import (
	"fmt"
)

const (
	DeepSeekBaseURL = "https://api.deepseek.com/v1"
)

type DeepSeekProvider struct{}

func init() {
	Register(&DeepSeekProvider{})
}

func (p *DeepSeekProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderDeepSeek,
		DisplayName: "DeepSeek",
		Description: "DeepSeek chat & reasoning models",
		DefaultURLs: map[ModelType]string{
			ModelTypeChat: DeepSeekBaseURL,
		},
		ModelTypes: []ModelType{
			ModelTypeChat,
		},
		RequiresAuth: true,
	}
}

func (p *DeepSeekProvider) ValidateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required for DeepSeek provider")
	}
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}
