package engine

var GlobalFactory *EngineFactory

type EngineFactory struct {
	llmProvider       LLMProvider
	embeddingProvider EmbeddingProvider
}

func NewEngineFactory() *EngineFactory {
	f := &EngineFactory{}
	GlobalFactory = f
	return f
}

func (f *EngineFactory) SetLLMProvider(p LLMProvider) {
	f.llmProvider = p
}

func (f *EngineFactory) SetEmbeddingProvider(p EmbeddingProvider) {
	f.embeddingProvider = p
}

func (f *EngineFactory) GetLLM() LLMProvider {
	return f.llmProvider
}

func (f *EngineFactory) GetEmbedding() EmbeddingProvider {
	return f.embeddingProvider
}

func GetGlobalFactory() *EngineFactory {
	return GlobalFactory
}
