package db

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/inaryzen/priotasks/models"
)

func setupTestDB(t *testing.T) *DbSQLite {
	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	tmpfile.Close()

	// Initialize test database
	db := NewDbSQLite()
	db.Init(tmpfile.Name())

	// Return the database and cleanup function
	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpfile.Name())
	})

	return db
}

func TestAddTagToTask_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create a test task
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	// Create and save a test tag
	tagId := uuid.New().String()
	if err := db.SaveTag(tagId); err != nil {
		t.Fatalf("failed to create test tag: %v", err)
	}

	// Test adding tag to task
	err := db.AddTagToTask(task.Id, tagId)
	if err != nil {
		t.Errorf("AddTagToTask failed: %v", err)
	}

	// Verify the tag was added by querying the database directly
	rows, err := db.instance.Query("SELECT task_id, tag_id FROM TasksTags WHERE task_id = ? AND tag_id = ?",
		task.Id, tagId)
	if err != nil {
		t.Fatalf("failed to query TasksTags: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Error("expected to find tag association in database, but found none")
	}
}

func TestAddTagToTask_InvalidTask(t *testing.T) {
	db := setupTestDB(t)

	// Create and save a test tag
	tagId := uuid.New().String()
	if err := db.SaveTag(tagId); err != nil {
		t.Fatalf("failed to create test tag: %v", err)
	}

	// Test adding invalid task/tag combination
	err := db.AddTagToTask("non-existent-task", tagId)
	if err == nil {
		t.Error("expected error when adding tag to non-existent task, got nil")
	}
}

func TestAddTagToTask_InvalidTag(t *testing.T) {
	db := setupTestDB(t)

	// Create a test task
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	// Test adding with non-existent tag
	err := db.AddTagToTask(task.Id, "non-existent-tag")
	if err == nil {
		t.Error("expected error when adding non-existent tag, got nil")
	}
}

func TestTaskTags_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create a test task
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	// Create and add multiple tags
	expectedTags := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for _, tagId := range expectedTags {
		if err := db.SaveTag(tagId); err != nil {
			t.Fatalf("failed to create tag: %v", err)
		}
		if err := db.AddTagToTask(task.Id, tagId); err != nil {
			t.Fatalf("failed to add tag to task: %v", err)
		}
	}

	// Test retrieving tags
	tags, err := db.TaskTags(task.Id)
	if err != nil {
		t.Fatalf("TaskTags failed: %v", err)
	}

	// Verify we got all expected tags
	if len(tags) != len(expectedTags) {
		t.Errorf("expected %d tags, got %d", len(expectedTags), len(tags))
	}

	// Verify each expected tag is present
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[string(tag)] = true
	}

	for _, expectedTag := range expectedTags {
		if !tagMap[expectedTag] {
			t.Errorf("expected to find tag %s, but it was missing", expectedTag)
		}
	}
}

func TestTaskTags_NonExistentTask(t *testing.T) {
	db := setupTestDB(t)

	// Test retrieving tags for non-existent task
	tags, err := db.TaskTags("non-existent-task")
	if err != nil {
		t.Fatalf("TaskTags failed unexpectedly: %v", err)
	}

	// Should return empty slice for non-existent task
	if len(tags) != 0 {
		t.Errorf("expected 0 tags for non-existent task, got %d", len(tags))
	}
}

func TestTags_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create multiple test tags
	expectedTags := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	// Save tags to database
	for _, tagId := range expectedTags {
		if err := db.SaveTag(tagId); err != nil {
			t.Fatalf("failed to create tag: %v", err)
		}
	}

	// Test retrieving all tags
	tags, err := db.Tags()
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	// Verify we got all expected tags
	if len(tags) != len(expectedTags) {
		t.Errorf("expected %d tags, got %d", len(expectedTags), len(tags))
	}

	// Create map of expected tags for easier lookup
	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[string(tag)] = true
	}

	// Verify each expected tag is present
	for _, expectedTag := range expectedTags {
		if !tagMap[expectedTag] {
			t.Errorf("expected to find tag %s, but it was missing", expectedTag)
		}
	}
}

func TestTags_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)

	// Test retrieving tags from empty database
	tags, err := db.Tags()
	if err != nil {
		t.Fatalf("Tags failed unexpectedly: %v", err)
	}

	// Should return empty slice when no tags exist
	if len(tags) != 0 {
		t.Errorf("expected 0 tags in empty database, got %d", len(tags))
	}
}
