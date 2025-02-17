package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupConfig struct {
	SourceDB     string
	BackupDir    string
	MaxBackups   int
	BackupFormat string
}

func NewBackupConfig() *BackupConfig {
	return &BackupConfig{
		SourceDB:     "./words.db",
		BackupDir:    "./backups",
		MaxBackups:   7,
		BackupFormat: "2006-01-02-150405", // YYYY-MM-DD-HHMMSS
	}
}

func CreateBackup(cfg *BackupConfig) error {
	// Ensure backup directory exists
	if err := os.MkdirAll(cfg.BackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup filename with timestamp
	backupFile := filepath.Join(
		cfg.BackupDir,
		fmt.Sprintf("backup-%s.db", time.Now().Format(cfg.BackupFormat)),
	)

	// Copy database file
	if err := copyFile(cfg.SourceDB, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Cleanup old backups
	return cleanupOldBackups(cfg)
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func cleanupOldBackups(cfg *BackupConfig) error {
	files, err := filepath.Glob(filepath.Join(cfg.BackupDir, "backup-*.db"))
	if err != nil {
		return err
	}

	if len(files) <= cfg.MaxBackups {
		return nil
	}

	// Sort files by modification time
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	fileInfos := make([]fileInfo, 0, len(files))
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{f, info.ModTime()})
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(fileInfos)-cfg.MaxBackups; i++ {
		if err := os.Remove(fileInfos[i].path); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", fileInfos[i].path, err)
		}
	}

	return nil
}
