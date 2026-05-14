package svcknowledge

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/ragflow/internal/dto/dtoknowledge"
)

type KnowledgeBaseSvc interface {
	Create(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseCreateReq) (*dtoknowledge.KnowledgeBaseCreateResp, error)
	Delete(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseDeleteReq) error
	Update(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseUpdateReq) error
	Detail(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseDetailReq) (*dtoknowledge.KnowledgeBaseDetailResp, error)
	List(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseListReq) (*dtoknowledge.KnowledgeBaseListResp, error)
}

type knowledgeBaseSvc struct {
}

var _ KnowledgeBaseSvc = (*knowledgeBaseSvc)(nil)

func NewKnowledgeBaseSvc() KnowledgeBaseSvc {
	return &knowledgeBaseSvc{}
}

func (svc *knowledgeBaseSvc) Create(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseCreateReq) (*dtoknowledge.KnowledgeBaseCreateResp, error) {
	return nil, nil
}

func (svc *knowledgeBaseSvc) Delete(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseDeleteReq) error {
	return nil
}

func (svc *knowledgeBaseSvc) Update(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseUpdateReq) error {
	return nil
}

func (svc *knowledgeBaseSvc) Detail(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseDetailReq) (*dtoknowledge.KnowledgeBaseDetailResp, error) {
	return nil, nil
}

func (svc *knowledgeBaseSvc) List(ctx *gin.Context, req *dtoknowledge.KnowledgeBaseListReq) (*dtoknowledge.KnowledgeBaseListResp, error) {
	return nil, nil
}