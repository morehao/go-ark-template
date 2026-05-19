package storage

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createMultipartFileHeader(t *testing.T, filename, content string) *multipart.FileHeader {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile failed: %v", err)
	}
	if _, err := part.Write([]byte(content)); err != nil {
		t.Fatalf("part.Write failed: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close failed: %v", err)
	}

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		t.Fatalf("ParseMultipartForm failed: %v", err)
	}
	return req.MultipartForm.File["file"][0]
}

func TestLocalFileService_SaveFile_GetFile_DeleteFile(t *testing.T) {
	basePath, err := os.MkdirTemp("", "storage-test")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer func() { _ = os.RemoveAll(basePath) }()

	svc := NewLocalFileService(basePath)
	ctx := context.Background()
	tenantID := uint(1)
	knowledgeID := "42"
	content := "hello world"

	fileHeader := createMultipartFileHeader(t, "test.txt", content)

	filePath, err := svc.SaveFile(ctx, fileHeader, tenantID, knowledgeID)
	if err != nil {
		t.Fatalf("SaveFile failed: %v", err)
	}
	if !strings.HasPrefix(filePath, localScheme) {
		t.Errorf("filePath should start with %s, got %s", localScheme, filePath)
	}

	reader, err := svc.GetFile(ctx, filePath)
	if err != nil {
		t.Fatalf("GetFile failed: %v", err)
	}
	defer func() { _ = reader.Close() }()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		t.Fatalf("ReadFrom failed: %v", err)
	}
	if buf.String() != content {
		t.Errorf("content = %q, want %q", buf.String(), content)
	}

	url, err := svc.GetFileURL(ctx, filePath)
	if err != nil {
		t.Fatalf("GetFileURL failed: %v", err)
	}
	if url != filePath {
		t.Errorf("GetFileURL = %q, want %q", url, filePath)
	}

	if err := svc.DeleteFile(ctx, filePath); err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}
	if _, err := svc.GetFile(ctx, filePath); err == nil {
		t.Error("expected error for deleted file, got nil")
	}
}

func TestLocalFileService_SaveBytes(t *testing.T) {
	basePath, err := os.MkdirTemp("", "storage-test")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer func() { _ = os.RemoveAll(basePath) }()

	svc := NewLocalFileService(basePath)
	ctx := context.Background()
	tenantID := uint(2)
	content := []byte("test data")

	filePath, err := svc.SaveBytes(ctx, content, tenantID, "test.txt", false)
	if err != nil {
		t.Fatalf("SaveBytes failed: %v", err)
	}
	if !strings.HasPrefix(filePath, localScheme) {
		t.Errorf("filePath should start with %s, got %s", localScheme, filePath)
	}

	trimmed := strings.TrimPrefix(filePath, localScheme)
	fullPath := filepath.Join(basePath, filepath.FromSlash(trimmed))
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("content = %q, want %q", string(data), string(content))
	}
}

func TestLocalFileService_CheckConnectivity(t *testing.T) {
	basePath, err := os.MkdirTemp("", "storage-test")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer func() { _ = os.RemoveAll(basePath) }()

	svc := NewLocalFileService(basePath)
	ctx := context.Background()

	if err := svc.CheckConnectivity(ctx); err != nil {
		t.Errorf("CheckConnectivity failed for existing dir: %v", err)
	}

	svcNonExistent := NewLocalFileService("/nonexistent/path")
	if err := svcNonExistent.CheckConnectivity(ctx); err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}

func TestLocalFileService_PathTraversalPrevention(t *testing.T) {
	basePath, err := os.MkdirTemp("", "storage-test")
	if err != nil {
		t.Fatalf("MkdirTemp failed: %v", err)
	}
	defer func() { _ = os.RemoveAll(basePath) }()

	svc := NewLocalFileService(basePath)
	ctx := context.Background()

	_, err = svc.SaveBytes(ctx, []byte("data"), 1, "../../etc/passwd", false)
	if err == nil {
		t.Error("expected error for path traversal in SaveBytes, got nil")
	}

	_, err = svc.GetFile(ctx, "local://../../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal in GetFile, got nil")
	}

	err = svc.DeleteFile(ctx, "local://../../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal in DeleteFile, got nil")
	}
}

func TestNewFileServiceFromStorageConfig(t *testing.T) {
	cfg := &EngineConfig{
		DefaultProvider: "local",
		Local: &LocalConfig{
			PathPrefix: "",
		},
	}

	svc, provider, err := NewFileServiceFromStorageConfig("", cfg, "/tmp")
	if err != nil {
		t.Fatalf("NewFileServiceFromStorageConfig failed: %v", err)
	}
	if provider != "local" {
		t.Errorf("provider = %q, want %q", provider, "local")
	}
	if svc == nil {
		t.Fatal("svc should not be nil")
	}

	_, _, err = NewFileServiceFromStorageConfig("invalid", cfg, "/tmp")
	if err == nil {
		t.Error("expected error for unsupported provider, got nil")
	}

	_, _, err = NewFileServiceFromStorageConfig("", &EngineConfig{}, "")
	if err == nil {
		t.Error("expected error for empty provider, got nil")
	}
}
