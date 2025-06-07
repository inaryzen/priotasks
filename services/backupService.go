package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/inaryzen/priotasks/common"
)

const (
	maxBackupFiles = 2
)

// BackupService handles SQLite database backup operations
type BackupService struct {
	baseDir string
}

// NewBackupService creates a new backup service instance
func NewBackupService(baseDir string) (*BackupService, error) {
	if baseDir == "" {
		return nil, fmt.Errorf("base directory cannot be empty")
	}
	return &BackupService{baseDir: baseDir}, nil
}

// CreateBackup makes a copy of the current database file
func (s *BackupService) CreateBackup() error {
	// Get source database path
	dbPath := filepath.Join(s.baseDir, "db.sqlite")

	// Generate backup file name with timestamp
	backupName := fmt.Sprintf("priotasks_db_backup_%s.db", time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(s.baseDir, backupName)

	// Copy the database file
	err := s.copyFile(dbPath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Clean up old backups
	return s.cleanupOldBackups()
}

// cleanupOldBackups ensures only the most recent backup files are kept
func (s *BackupService) cleanupOldBackups() error {
	pattern := filepath.Join(s.baseDir, "priotasks_db_backup_*.db")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list backup files: %w", err)
	}

	if len(matches) <= maxBackupFiles {
		return nil
	}

	// Sort modification time DESC
	sort.Slice(matches, func(i, j int) bool {
		iInfo, _ := os.Stat(matches[i])
		jInfo, _ := os.Stat(matches[j])
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	// Remove older backups
	for _, file := range matches[maxBackupFiles:] {
		err := os.Remove(file)
		if err != nil {
			common.Debug("Failed to remove old backup %s: %v", file, err)
		} else {
			log.Printf("removed backup: %v\n", file)
		}
	}

	return nil
}

// copyFile creates a copy of a file with proper permissions
func (s *BackupService) copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err == nil {
		log.Printf("created backup: %v\n", dst)
	}
	return err
}
