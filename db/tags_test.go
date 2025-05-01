package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/inaryzen/priotasks/models"
)

func TestSaveAndFindTag_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create and save a task first
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	tag := "test-tag"
	err := db.SaveTag(tag)
	if err != nil {
		t.Fatalf("SaveTag failed: %v", err)
	}

	// Add tag to task
	err = db.AddTagToTask(task.Id, tag)
	if err != nil {
		t.Fatalf("AddTagToTask failed: %v", err)
	}

	tags, err := db.Tags()
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
	if string(tags[0]) != tag {
		t.Errorf("expected tag %s, got %s", tag, tags[0])
	}
}

func TestAddAndGetTaskTags_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create and save a task first
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	// Save unique tags first
	tags := []string{"tag1", "tag2", "tag3"}
	uniqueTags := make(map[string]bool)
	for _, tag := range tags {
		if !uniqueTags[tag] {
			err := db.SaveTag(tag)
			if err != nil {
				t.Fatalf("SaveTag failed: %v", err)
			}
			uniqueTags[tag] = true
		}
	}

	// Add tags to task
	for _, tag := range tags {
		err := db.AddTagToTask(task.Id, tag)
		if err != nil {
			t.Fatalf("AddTagToTask failed: %v", err)
		}
	}

	// Get tags for task
	taskTags, err := db.TaskTags(task.Id)
	if err != nil {
		t.Fatalf("TaskTags failed: %v", err)
	}

	if len(taskTags) != len(tags) {
		t.Errorf("expected %d tags, got %d", len(tags), len(taskTags))
	}

	// Verify all tags are present
	for _, tag := range tags {
		found := false
		for _, taskTag := range taskTags {
			if string(taskTag) == tag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("tag %s not found in task tags", tag)
		}
	}
}

func TestDeleteTagFromTask_Success(t *testing.T) {
	db := setupTestDB(t)

	// Create and save a task first
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	tag := "test-tag"

	// Setup
	err := db.SaveTag(tag)
	if err != nil {
		t.Fatalf("SaveTag failed: %v", err)
	}
	err = db.AddTagToTask(task.Id, tag)
	if err != nil {
		t.Fatalf("AddTagToTask failed: %v", err)
	}

	// Delete tag
	err = db.DeleteTagFromTask(task.Id, tag)
	if err != nil {
		t.Fatalf("DeleteTagFromTask failed: %v", err)
	}

	// Verify tag was deleted
	taskTags, err := db.TaskTags(task.Id)
	if err != nil {
		t.Fatalf("TaskTags failed: %v", err)
	}

	if len(taskTags) != 0 {
		t.Errorf("expected 0 tags after deletion, got %d", len(taskTags))
	}
}

func TestDeleteTagFromTask_NonExistent(t *testing.T) {
	db := setupTestDB(t)

	// Create and save a task first
	task := models.Task{
		Id:    uuid.New().String(),
		Title: "Test Task",
	}
	if err := db.SaveTask(task); err != nil {
		t.Fatalf("failed to create test task: %v", err)
	}

	err := db.DeleteTagFromTask(task.Id, "non-existent-tag")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for non-existent tag, got: %v", err)
	}
}

func TestTasksTags_Success(t *testing.T) {
	db := setupTestDB(t)

	taskIds := []string{"task1", "task2", "task3"}
	tagsByTask := map[string][]string{
		"task1": {"tag1", "tag2"},
		"task2": {"tag2", "tag3"},
		"task3": {"tag1", "tag3"},
	}

	// Create tasks first
	for _, taskId := range taskIds {
		task := models.Task{
			Id:    taskId,
			Title: "Test Task " + taskId,
		}
		if err := db.SaveTask(task); err != nil {
			t.Fatalf("failed to create task %s: %v", taskId, err)
		}
	}

	// Save unique tags first
	uniqueTags := make(map[string]bool)
	for _, tags := range tagsByTask {
		for _, tag := range tags {
			if !uniqueTags[tag] {
				if err := db.SaveTag(tag); err != nil {
					t.Fatalf("failed to save tag %s: %v", tag, err)
				}
				uniqueTags[tag] = true
			}
		}
	}

	// Create task-tag associations
	for taskId, tags := range tagsByTask {
		for _, tag := range tags {
			if err := db.AddTagToTask(taskId, tag); err != nil {
				t.Fatalf("failed to add tag %s to task %s: %v", tag, taskId, err)
			}
		}
	}

	// Get tags for all tasks
	result, err := db.TasksTags(taskIds)
	if err != nil {
		t.Fatalf("TasksTags failed: %v", err)
	}

	// Verify results
	for taskId, expectedTags := range tagsByTask {
		taskTags := result[taskId]
		if len(taskTags) != len(expectedTags) {
			t.Errorf("task %s: expected %d tags, got %d", taskId, len(expectedTags), len(taskTags))
		}

		for _, expectedTag := range expectedTags {
			found := false
			for _, taskTag := range taskTags {
				if string(taskTag) == expectedTag {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("task %s: tag %s not found in result", taskId, expectedTag)
			}
		}
	}
}

func TestTasksTags_EmptyInput(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.TasksTags([]string{})
	if err != nil {
		t.Fatalf("TasksTags failed with empty input: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result for empty input, got %d entries", len(result))
	}
}

func TestTasksTags_NonExistentTasks(t *testing.T) {
	db := setupTestDB(t)

	result, err := db.TasksTags([]string{"non-existent-task"})
	if err != nil {
		t.Fatalf("TasksTags failed with non-existent task: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result for non-existent task, got %d entries", len(result))
	}
}
