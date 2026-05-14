package dtoknowledge

type KnowledgeBaseCreateReq struct {
	Name           string `json:"name" validate:"required" label:"知识库名称"`
	Description    string `json:"description" label:"知识库描述"`
	EmbeddingModel string `json:"embeddingModel" validate:"required" label:"embedding模型"`
	VectorStoreType string `json:"vectorStoreType" validate:"required" label:"向量库类型"`
	PermissionType  string `json:"permissionType" validate:"required" label:"权限类型"`
	ChunkMethod     string `json:"chunkMethod" validate:"required" label:"分块方式"`
}

type KnowledgeBaseUpdateReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"`
	KnowledgeBaseCreateReq
}

type KnowledgeBaseDetailReq struct {
	ID uint `json:"id" form:"id" validate:"required" label:"数据自增id"`
}

type KnowledgeBaseListReq struct {
	Page     int    `form:"page" validate:"required" label:"页码"`
	PageSize int    `form:"pageSize" validate:"required" label:"每页数量"`
	Name     string `form:"name" label:"知识库名称"`
	Status   string `form:"status" label:"状态"`
}

type KnowledgeBaseDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"`
}

type DocumentUploadReq struct {
	KnowledgeBaseID uint   `json:"knowledgeBaseID" validate:"required" label:"知识库id"`
	Name            string `json:"name" validate:"required" label:"文档名称"`
	Type            string `json:"type" validate:"required" label:"文档类型"`
	Location        string `json:"location" validate:"required" label:"文档位置"`
}

type DocumentListReq struct {
	Page            int    `form:"page" validate:"required" label:"页码"`
	PageSize        int    `form:"pageSize" validate:"required" label:"每页数量"`
	KnowledgeBaseID uint   `form:"knowledgeBaseID" label:"知识库id"`
	Name            string `form:"name" label:"文档名称"`
}

type DocumentDeleteReq struct {
	ID uint `json:"id" validate:"required" label:"数据自增id"`
}