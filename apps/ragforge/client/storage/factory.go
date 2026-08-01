package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NewFileServiceFromStorageConfig(
	provider string,
	cfg *EngineConfig,
	localBaseDir string,
) (FileService, string, error) {
	p := strings.ToLower(strings.TrimSpace(provider))
	if p == "" && cfg != nil {
		p = strings.ToLower(strings.TrimSpace(cfg.DefaultProvider))
	}
	if p == "" {
		return nil, "", fmt.Errorf("empty provider")
	}

	if localBaseDir == "" {
		localBaseDir = strings.TrimSpace(os.Getenv("LOCAL_STORAGE_BASE_DIR"))
	}
	if localBaseDir == "" {
		localBaseDir = "./data/files"
	}

	switch p {
	case "local":
		baseDir := localBaseDir
		if cfg != nil && cfg.Local != nil {
			prefix := strings.TrimSpace(cfg.Local.PathPrefix)
			prefix = strings.Trim(prefix, "/\\")
			if prefix != "" {
				baseDir = filepath.Join(baseDir, prefix)
			}
		}
		return NewLocalFileService(baseDir), p, nil

	default:
		return nil, p, fmt.Errorf("unsupported provider %q", p)
	}
}
