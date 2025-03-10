package db

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/inaryzen/priotasks/models"
)

func TestSaveAndFindSettings_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create test settings
	settings := models.Settings{
		Id: uuid.New().String(),
		TasksQuery: models.TasksQuery{
			FilterCompleted:   true,
			FilterIncompleted: false,
			CompletedFrom:     time.Now().AddDate(0, 0, -7),
			CompletedTo:       time.Now(),
			SortColumn:        models.Priority,
			SortDirection:     models.Desc,
			FilterWip:         true,
			FilterNonWip:      false,
			Planned:           true,
			NonPlanned:        false,
			Tags:              []models.TaskTag{"tag1", "tag2"},
		},
	}

	// Save settings
	err := db.SaveSettings(settings)
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	// Find settings
	found, err := db.FindSettings(settings.Id)
	if err != nil {
		t.Fatalf("FindSettings failed: %v", err)
	}

	// Verify fields match
	if found.Id != settings.Id {
		t.Errorf("expected settings ID %s, got %s", settings.Id, found.Id)
	}
	if found.TasksQuery.FilterCompleted != settings.TasksQuery.FilterCompleted {
		t.Errorf("expected FilterCompleted %v, got %v",
			settings.TasksQuery.FilterCompleted, found.TasksQuery.FilterCompleted)
	}
	if found.TasksQuery.FilterWip != settings.TasksQuery.FilterWip {
		t.Errorf("expected FilterWip %v, got %v",
			settings.TasksQuery.FilterWip, found.TasksQuery.FilterWip)
	}
	if found.TasksQuery.SortColumn != settings.TasksQuery.SortColumn {
		t.Errorf("expected SortColumn %v, got %v",
			settings.TasksQuery.SortColumn, found.TasksQuery.SortColumn)
	}
	if len(found.TasksQuery.Tags) != len(settings.TasksQuery.Tags) {
		t.Errorf("expected %d tags, got %d",
			len(settings.TasksQuery.Tags), len(found.TasksQuery.Tags))
	}
}

func TestFindSettings_NonExistent(t *testing.T) {
	db := setupTestDB(t)

	_, err := db.FindSettings("non-existent-id")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for non-existent settings, got: %v", err)
	}
}

func TestSaveSettings_Update(t *testing.T) {
	db := setupTestDB(t)

	// Initial settings
	settings := models.Settings{
		Id: uuid.New().String(),
		TasksQuery: models.TasksQuery{
			FilterCompleted: true,
			SortColumn:      models.Priority,
			Tags:            []models.TaskTag{"tag1"},
		},
	}

	// Save initial settings
	err := db.SaveSettings(settings)
	if err != nil {
		t.Fatalf("Initial SaveSettings failed: %v", err)
	}

	// Modify settings
	settings.TasksQuery.FilterCompleted = false
	settings.TasksQuery.Tags = []models.TaskTag{"tag1", "tag2"}

	// Update settings
	err = db.SaveSettings(settings)
	if err != nil {
		t.Fatalf("Update SaveSettings failed: %v", err)
	}

	// Verify updates
	found, err := db.FindSettings(settings.Id)
	if err != nil {
		t.Fatalf("FindSettings failed after update: %v", err)
	}

	if found.TasksQuery.FilterCompleted != false {
		t.Error("FilterCompleted was not updated")
	}
	if len(found.TasksQuery.Tags) != 2 {
		t.Errorf("expected 2 tags after update, got %d", len(found.TasksQuery.Tags))
	}
}

func TestSaveSettings_WithDates(t *testing.T) {
	db := setupTestDB(t)

	completedFrom := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	completedTo := completedFrom.AddDate(0, 0, 15)

	settings := models.Settings{
		Id: uuid.New().String(),
		TasksQuery: models.TasksQuery{
			CompletedFrom: completedFrom,
			CompletedTo:   completedTo,
		},
	}

	// Save settings
	err := db.SaveSettings(settings)
	if err != nil {
		t.Fatalf("SaveSettings failed: %v", err)
	}

	// Find settings
	found, err := db.FindSettings(settings.Id)
	if err != nil {
		t.Fatalf("FindSettings failed: %v", err)
	}

	// Verify dates
	if !found.TasksQuery.CompletedFrom.Equal(completedFrom) {
		t.Errorf("CompletedFrom date mismatch: expected %v, got %v",
			completedFrom, found.TasksQuery.CompletedFrom)
	}
	if !found.TasksQuery.CompletedTo.Equal(completedTo) {
		t.Errorf("CompletedTo date mismatch: expected %v, got %v",
			completedTo, found.TasksQuery.CompletedTo)
	}
}
