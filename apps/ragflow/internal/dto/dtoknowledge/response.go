package dtoknowledge

type KnowledgeBaseCreateResp struct {
	ID uint `json:"id"`
}

type KnowledgeBaseDetailResp struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	EmbeddingModel string `json:"embeddingModel"`
	VectorStoreType string `json:"vectorStoreType"`
	PermissionType  string `json:"permissionType"`
	Status        string `json:"status"`
	ChunkMethod   string `json:"chunkMethod"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

type KnowledgeBaseListItem struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	EmbeddingModel string `json:"embeddingModel"`
	VectorStoreType string `json:"vectorStoreType"`
	PermissionType  string `json:"permissionType"`
	Status        string `json:"status"`
	ChunkMethod   string `json:"chunkMethod"`
	UpdatedAt     int64  `json:"updatedAt"`
}

type KnowledgeBaseListResp struct {
	List  []KnowledgeBaseListItem `json:"list"`
	Total int64                   `json:"total"`
}

type DocumentUploadResp struct {
	ID uint `json:"id"`
}

type DocumentListItem struct {
	ID               uint   `json:"id"`
	KnowledgeBaseID  uint   `json:"knowledgeBaseID"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Location         string `json:"location"`
	Size             int64  `json:"size"`
	Status           string `json:"status"`
	ChunkStatus      string `json:"chunkStatus"`
	VectorStatus     string `json:"vectorStatus"`
	UpdatedAt        int64  `json:"updatedAt"`
}

type DocumentListResp struct {
	List  []DocumentListItem `json:"list"`
	Total int64              `json:"total"`
}