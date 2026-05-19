package provider

import (
	"fmt"
)

type GenericProvider struct{}

func init() {
	Register(&GenericProvider{})
}

func (p *GenericProvider) Info() ProviderInfo {
	return ProviderInfo{
		Name:        ProviderGeneric,
		DisplayName: "通用 OpenAI 兼容",
		Description: "兼容 OpenAI API 格式的自定义部署服务",
		DefaultURLs: map[ModelType]string{
			ModelTypeChat:      "",
			ModelTypeEmbedding: "",
			ModelTypeRerank:    "",
		},
		ModelTypes: []ModelType{
			ModelTypeChat,
			ModelTypeEmbedding,
			ModelTypeRerank,
		},
		RequiresAuth: false,
	}
}

func (p *GenericProvider) ValidateConfig(config *Config) error {
	if config.ModelName == "" {
		return fmt.Errorf("model name is required")
	}
	return nil
}
