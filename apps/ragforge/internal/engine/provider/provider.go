package provider

import (
	"fmt"
	"strings"
	"sync"
)

type ModelType string

const (
	ModelTypeChat      ModelType = "chat"
	ModelTypeEmbedding ModelType = "embedding"
	ModelTypeRerank    ModelType = "rerank"
)

type ProviderName string

const (
	ProviderOpenAI      ProviderName = "openai"
	ProviderAnthropic   ProviderName = "anthropic"
	ProviderAliyun      ProviderName = "aliyun"
	ProviderZhipu       ProviderName = "zhipu"
	ProviderOpenRouter  ProviderName = "openrouter"
	ProviderSiliconFlow ProviderName = "siliconflow"
	ProviderJina        ProviderName = "jina"
	ProviderGeneric     ProviderName = "generic"
	ProviderDeepSeek    ProviderName = "deepseek"
	ProviderGemini      ProviderName = "gemini"
	ProviderVolcengine  ProviderName = "volcengine"
	ProviderHunyuan     ProviderName = "hunyuan"
	ProviderMiniMax     ProviderName = "minimax"
	ProviderMimo        ProviderName = "mimo"
	ProviderGPUStack    ProviderName = "gpustack"
	ProviderMoonshot    ProviderName = "moonshot"
	ProviderModelScope  ProviderName = "modelscope"
	ProviderQianfan     ProviderName = "qianfan"
	ProviderQiniu       ProviderName = "qiniu"
	ProviderLongCat     ProviderName = "longcat"
	ProviderLKEAP       ProviderName = "lkeap"
	ProviderNvidia      ProviderName = "nvidia"
	ProviderNovita      ProviderName = "novita"
	ProviderAzureOpenAI ProviderName = "azure_openai"
)

type ProviderInfo struct {
	Name         ProviderName
	DisplayName  string
	Description  string
	DefaultURLs  map[ModelType]string
	ModelTypes   []ModelType
	RequiresAuth bool
	ExtraFields  []ExtraFieldConfig
}

func (p ProviderInfo) GetDefaultURL(modelType ModelType) string {
	if url, ok := p.DefaultURLs[modelType]; ok {
		return url
	}
	if url, ok := p.DefaultURLs[ModelTypeChat]; ok {
		return url
	}
	return ""
}

type ExtraFieldConfig struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Default     string `json:"default"`
	Placeholder string `json:"placeholder"`
	Options     []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"options,omitempty"`
}

type Config struct {
	Provider  ProviderName   `json:"provider"`
	BaseURL   string         `json:"base_url"`
	APIKey    string         `json:"api_key"`
	ModelName string         `json:"model_name"`
	ModelID   string         `json:"model_id"`
	Extra     map[string]any `json:"extra,omitempty"`
}

type Provider interface {
	Info() ProviderInfo
	ValidateConfig(config *Config) error
}

var (
	registryMu sync.RWMutex
	registry   = make(map[ProviderName]Provider)
)

func Register(p Provider) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[p.Info().Name] = p
}

func Get(name ProviderName) (Provider, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	p, ok := registry[name]
	return p, ok
}

func GetOrDefault(name ProviderName) Provider {
	p, ok := Get(name)
	if ok {
		return p
	}
	p, _ = Get(ProviderGeneric)
	return p
}

func List() []ProviderInfo {
	registryMu.RLock()
	defer registryMu.RUnlock()

	result := make([]ProviderInfo, 0, len(registry))
	for _, p := range registry {
		result = append(result, p.Info())
	}
	return result
}

func ListByModelType(modelType ModelType) []ProviderInfo {
	registryMu.RLock()
	defer registryMu.RUnlock()

	result := make([]ProviderInfo, 0)
	for _, p := range registry {
		info := p.Info()
		for _, t := range info.ModelTypes {
			if t == modelType {
				result = append(result, info)
				break
			}
		}
	}
	return result
}

func DetectProvider(baseURL string) ProviderName {
	switch {
	case containsAny(baseURL, "dashscope.aliyuncs.com"):
		return ProviderAliyun
	case containsAny(baseURL, "open.bigmodel.cn", "zhipu"):
		return ProviderZhipu
	case containsAny(baseURL, "openrouter.ai"):
		return ProviderOpenRouter
	case containsAny(baseURL, "siliconflow.cn"):
		return ProviderSiliconFlow
	case containsAny(baseURL, "api.jina.ai"):
		return ProviderJina
	case containsAny(baseURL, "openai.azure.com"):
		return ProviderAzureOpenAI
	case containsAny(baseURL, "api.openai.com"):
		return ProviderOpenAI
	case containsAny(baseURL, "api.anthropic.com"):
		return ProviderAnthropic
	case containsAny(baseURL, "api.deepseek.com"):
		return ProviderDeepSeek
	case containsAny(baseURL, "generativelanguage.googleapis.com"):
		return ProviderGemini
	case containsAny(baseURL, "volces.com", "volcengine"):
		return ProviderVolcengine
	case containsAny(baseURL, "hunyuan.cloud.tencent.com"):
		return ProviderHunyuan
	case containsAny(baseURL, "minimax.io", "minimaxi.com"):
		return ProviderMiniMax
	case containsAny(baseURL, "xiaomimimo.com"):
		return ProviderMimo
	case containsAny(baseURL, "gpustack"):
		return ProviderGPUStack
	case containsAny(baseURL, "modelscope.cn"):
		return ProviderModelScope
	case containsAny(baseURL, "qiniuapi.com", "qiniu"):
		return ProviderQiniu
	case containsAny(baseURL, "moonshot.ai"):
		return ProviderMoonshot
	case containsAny(baseURL, "qianfan.baidubce.com", "baidubce.com"):
		return ProviderQianfan
	case containsAny(baseURL, "longcat.chat"):
		return ProviderLongCat
	case containsAny(baseURL, "lkeap.cloud.tencent.com", "api.lkeap"):
		return ProviderLKEAP
	case containsAny(baseURL, "nvidia.com"):
		return ProviderNvidia
	case containsAny(baseURL, "api.novita.ai", "novita.ai"):
		return ProviderNovita
	default:
		return ProviderGeneric
	}
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

type ModelParams struct {
	Name     string
	ID       string
	BaseURL  string
	APIKey   string
	Provider string
}

func NewConfigFromModel(model *ModelParams) (*Config, error) {
	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	providerName := ProviderName(model.Provider)
	if providerName == "" {
		providerName = DetectProvider(model.BaseURL)
	}

	return &Config{
		Provider:  providerName,
		BaseURL:   model.BaseURL,
		APIKey:    model.APIKey,
		ModelName: model.Name,
		ModelID:   model.ID,
	}, nil
}
