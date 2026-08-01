package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/morehao/goark/pkg/dbclient"
	"github.com/morehao/golib/glog"
	"gorm.io/gorm"
)

func RunMigrations(ctx context.Context) error {
	db := dbclient.RagForgeDB(ctx)
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	if err := ensureMigrationTable(db); err != nil {
		return fmt.Errorf("ensure migration table fail: %w", err)
	}

	migrationDir := "migrations"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		execPath, _ := os.Executable()
		migrationDir = filepath.Join(filepath.Dir(execPath), "migrations")
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			glog.Infof(ctx, "[database.RunMigrations] migration dir not found, skip")
			return nil
		}
	}

	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("read migration dir fail: %w", err)
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	for _, filename := range upFiles {
		applied, err := isMigrationApplied(db, filename)
		if err != nil {
			return err
		}
		if applied {
			glog.Infof(ctx, "[database.RunMigrations] skip already applied: %s", filename)
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationDir, filename))
		if err != nil {
			return fmt.Errorf("read migration %s fail: %w", filename, err)
		}

		glog.Infof(ctx, "[database.RunMigrations] applying: %s", filename)
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("apply migration %s fail: %w", filename, err)
		}

		if err := recordMigration(db, filename); err != nil {
			return fmt.Errorf("record migration %s fail: %w", filename, err)
		}
		glog.Infof(ctx, "[database.RunMigrations] applied: %s", filename)
	}

	return nil
}

func ensureMigrationTable(db *gorm.DB) error {
	raw := `CREATE TABLE IF NOT EXISTS rg_migration (
		id BIGSERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	return db.Exec(raw).Error
}

func isMigrationApplied(db *gorm.DB, filename string) (bool, error) {
	var count int64
	err := db.Table("rg_migration").Where("filename = ?", filename).Count(&count).Error
	return count > 0, err
}

func recordMigration(db *gorm.DB, filename string) error {
	return db.Exec("INSERT INTO rg_migration (filename) VALUES (?)", filename).Error
}
