package services

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

func TestBackupService_CreateBackup(t *testing.T) {
	// Setup temporary directory for test
	tempDir, err := os.MkdirTemp("", "backuptest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a fake DB file
	dbPath := filepath.Join(tempDir, "db.sqlite")
	content := []byte("test database content")
	if err := os.WriteFile(dbPath, content, 0600); err != nil {
		t.Fatalf("Failed to create test db file: %v", err)
	}

	// Create backup service with temp dir
	service := &BackupService{baseDir: tempDir}

	// Test creating multiple backups
	for i := 0; i < 3; i++ {
		if err := service.CreateBackup(); err != nil {
			t.Errorf("CreateBackup failed: %v", err)
		}
		time.Sleep(time.Second) // Ensure unique timestamps
	}

	// Check number of backup files
	pattern := filepath.Join(tempDir, "priotasks_db_backup_*.db")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Failed to list backups: %v", err)
	}

	if len(matches) != maxBackupFiles {
		t.Errorf("Expected %d backup files, got %d", maxBackupFiles, len(matches))
	}

	// Verify backup content
	latest := matches[0]
	backupContent, err := os.ReadFile(latest)
	if err != nil {
		t.Fatalf("Failed to read backup: %v", err)
	}

	if string(backupContent) != string(content) {
		t.Error("Backup content does not match original")
	}
}

func TestExploreSliceSort(t *testing.T) {
	numbers := []int{5, 2, 8, 1, 9, 3}

	// Sort in ASC order using sort.Slice
	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i] < numbers[j]
	})

	for i := 1; i < len(numbers); i++ {
		if numbers[i-1] > numbers[i] {
			t.Errorf("Numbers not properly sorted at index %d: %v should be less than %v",
				i, numbers[i-1], numbers[i])
		}
	}

	t.Logf("Sorted numbers: %v", numbers)
}
