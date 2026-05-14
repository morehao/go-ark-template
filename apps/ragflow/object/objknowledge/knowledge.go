package objknowledge

type KnowledgeBaseBaseInfo struct {
	Name           string `json:"name" form:"name"`
	Description    string `json:"description" form:"description"`
	EmbeddingModel string `json:"embeddingModel" form:"embeddingModel"`
	VectorStoreType string `json:"vectorStoreType" form:"vectorStoreType"`
	PermissionType  string `json:"permissionType" form:"permissionType"`
	ChunkMethod     string `json:"chunkMethod" form:"chunkMethod"`
}