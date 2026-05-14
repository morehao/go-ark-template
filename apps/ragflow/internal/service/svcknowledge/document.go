package svcknowledge

import (
	"github.com/gin-gonic/gin"
	"github.com/morehao/goark/apps/ragflow/internal/dto/dtoknowledge"
)

type DocumentSvc interface {
	Upload(ctx *gin.Context, req *dtoknowledge.DocumentUploadReq) (*dtoknowledge.DocumentUploadResp, error)
	List(ctx *gin.Context, req *dtoknowledge.DocumentListReq) (*dtoknowledge.DocumentListResp, error)
	Delete(ctx *gin.Context, req *dtoknowledge.DocumentDeleteReq) error
}

type documentSvc struct {
}

var _ DocumentSvc = (*documentSvc)(nil)

func NewDocumentSvc() DocumentSvc {
	return &documentSvc{}
}

func (svc *documentSvc) Upload(ctx *gin.Context, req *dtoknowledge.DocumentUploadReq) (*dtoknowledge.DocumentUploadResp, error) {
	return nil, nil
}

func (svc *documentSvc) List(ctx *gin.Context, req *dtoknowledge.DocumentListReq) (*dtoknowledge.DocumentListResp, error) {
	return nil, nil
}

func (svc *documentSvc) Delete(ctx *gin.Context, req *dtoknowledge.DocumentDeleteReq) error {
	return nil
}