package storage

import (
	"context"
	"io"
	"mime/multipart"
)

type FileService interface {
	CheckConnectivity(ctx context.Context) error
	SaveFile(ctx context.Context, file *multipart.FileHeader, tenantID uint, knowledgeID string) (string, error)
	SaveBytes(ctx context.Context, data []byte, tenantID uint, fileName string, temp bool) (string, error)
	GetFile(ctx context.Context, filePath string) (io.ReadCloser, error)
	GetFileURL(ctx context.Context, filePath string) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
}
