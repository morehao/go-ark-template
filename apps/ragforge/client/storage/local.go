package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/morehao/golib/glog"
)

type localFileService struct {
	baseDir string
}

const (
	localScheme    = "local://"
	exportsDirName = "exports"
)

func (s *localFileService) CheckConnectivity(ctx context.Context) error {
	info, err := os.Stat(s.baseDir)
	if err != nil {
		return fmt.Errorf("storage directory not accessible: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("storage path is not a directory: %s", s.baseDir)
	}
	return nil
}

func NewLocalFileService(baseDir string) FileService {
	return &localFileService{
		baseDir: baseDir,
	}
}

func (s *localFileService) SaveFile(ctx context.Context,
	file *multipart.FileHeader, tenantID uint, knowledgeID string,
) (string, error) {
	glog.Infof(ctx, "[storage.local.SaveFile] name=%s, size=%d, tenantID=%d, knowledgeID=%s",
		file.Filename, file.Size, tenantID, knowledgeID)

	dir := filepath.Join(s.baseDir, fmt.Sprintf("%d", tenantID), knowledgeID)
	if _, err := safePathUnderBase(s.baseDir, dir); err != nil {
		glog.Errorf(ctx, "[storage.local.SaveFile] path traversal denied: %v", err)
		return "", fmt.Errorf("invalid path: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		glog.Errorf(ctx, "[storage.local.SaveFile] mkdir fail: %v", err)
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		glog.Errorf(ctx, "[storage.local.SaveFile] open source fail: %v", err)
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = src.Close() }()

	dst, err := os.Create(filePath)
	if err != nil {
		glog.Errorf(ctx, "[storage.local.SaveFile] create destination fail: %v", err)
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = dst.Close() }()

	if _, err := io.Copy(dst, src); err != nil {
		glog.Errorf(ctx, "[storage.local.SaveFile] copy fail: %v", err)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	relPath, _ := filepath.Rel(s.baseDir, filePath)
	return localScheme + filepath.ToSlash(relPath), nil
}

func (s *localFileService) SaveBytes(ctx context.Context, data []byte, tenantID uint, fileName string, temp bool) (string, error) {
	glog.Infof(ctx, "[storage.local.SaveBytes] fileName=%s, size=%d, tenantID=%d, temp=%v", fileName, len(data), tenantID, temp)

	safeName, err := safeFileName(fileName)
	if err != nil {
		glog.Errorf(ctx, "[storage.local.SaveBytes] invalid fileName: %v", err)
		return "", fmt.Errorf("invalid file name: %w", err)
	}

	dir := filepath.Join(s.baseDir, fmt.Sprintf("%d", tenantID), exportsDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		glog.Errorf(ctx, "[storage.local.SaveBytes] mkdir fail: %v", err)
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	ext := filepath.Ext(safeName)
	baseName := safeName[:len(safeName)-len(ext)]
	uniqueFileName := fmt.Sprintf("%s_%d%s", baseName, time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, uniqueFileName)

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		glog.Errorf(ctx, "[storage.local.SaveBytes] write file fail: %v", err)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	relPath, _ := filepath.Rel(s.baseDir, filePath)
	return localScheme + filepath.ToSlash(relPath), nil
}

func (s *localFileService) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	glog.Infof(ctx, "[storage.local.GetFile] path=%s", filePath)

	candidate := s.normalizePathForBase(filePath)
	resolved, err := safePathUnderBase(s.baseDir, candidate)
	if err != nil {
		glog.Errorf(ctx, "[storage.local.GetFile] path traversal denied: %v", err)
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	file, err := os.Open(resolved)
	if err != nil {
		glog.Errorf(ctx, "[storage.local.GetFile] open fail: %v", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (s *localFileService) DeleteFile(ctx context.Context, filePath string) error {
	glog.Infof(ctx, "[storage.local.DeleteFile] path=%s", filePath)

	candidate := s.normalizePathForBase(filePath)
	resolved, err := safePathUnderBase(s.baseDir, candidate)
	if err != nil {
		glog.Errorf(ctx, "[storage.local.DeleteFile] path traversal denied: %v", err)
		return fmt.Errorf("invalid file path: %w", err)
	}

	if err := os.Remove(resolved); err != nil {
		glog.Errorf(ctx, "[storage.local.DeleteFile] remove fail: %v", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *localFileService) GetFileURL(ctx context.Context, filePath string) (string, error) {
	normalized := filePath
	if !strings.HasPrefix(filePath, localScheme) {
		relPath, err := filepath.Rel(s.baseDir, filePath)
		if err != nil {
			normalized = filePath
		} else {
			normalized = localScheme + filepath.ToSlash(relPath)
		}
	}
	return normalized, nil
}

func (s *localFileService) normalizePathForBase(filePath string) string {
	if strings.HasPrefix(filePath, localScheme) {
		relPath := strings.TrimPrefix(filePath, localScheme)
		return filepath.Join(s.baseDir, filepath.FromSlash(relPath))
	}

	clean := filepath.Clean(strings.TrimSpace(filePath))
	if clean == "." || clean == "" {
		return clean
	}
	if filepath.IsAbs(clean) {
		return clean
	}

	baseClean := filepath.Clean(s.baseDir)
	baseNoSlash := strings.Trim(baseClean, string(filepath.Separator))
	cleanNoDot := strings.TrimPrefix(clean, "."+string(filepath.Separator))
	if strings.HasPrefix(cleanNoDot, baseNoSlash+string(filepath.Separator)) {
		cleanNoDot = strings.TrimPrefix(cleanNoDot, baseNoSlash+string(filepath.Separator))
	}
	return filepath.Join(baseClean, cleanNoDot)
}

func safeFileName(name string) (string, error) {
	clean := path.Clean(name)
	if strings.Contains(clean, "/") || strings.Contains(clean, "\\") || clean == "." || clean == ".." {
		return "", fmt.Errorf("unsafe file name: %s", name)
	}
	return clean, nil
}

func safePathUnderBase(baseDir, candidate string) (string, error) {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("cannot resolve base dir: %w", err)
	}
	absCandidate, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("cannot resolve candidate path: %w", err)
	}
	cleanBase := filepath.Clean(absBase)
	cleanCandidate := filepath.Clean(absCandidate)
	if !strings.HasPrefix(cleanCandidate, cleanBase+string(filepath.Separator)) && cleanCandidate != cleanBase {
		return "", fmt.Errorf("path traversal denied: %s is not under %s", candidate, baseDir)
	}
	return cleanCandidate, nil
}
