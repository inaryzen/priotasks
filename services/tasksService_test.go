package services

import (
	"fmt"
	"testing"

	"slices"

	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

// MockDB implements db.Db interface for testing
type MockDB struct {
	tasks      map[string]models.Task
	migrations map[string]bool
	tags       map[string]bool
	taskTags   map[string][]models.TaskTag
}

func (m *MockDB) DeleteTag(tagId string) error             { return nil }
func (m *MockDB) DeleteTagFromAllTasks(tagId string) error { return nil }

func (m *MockDB) Init(p string) {}
func (m *MockDB) Close()        {}
func (m *MockDB) Tasks() ([]models.Task, error) {
	var result []models.Task
	for _, value := range m.tasks {
		result = append(result, value)
	}
	return result, nil
}
func (m *MockDB) FindTasks(query models.TasksQuery) ([]models.Task, error) { return nil, nil }
func (m *MockDB) DeleteTask(taskId string) error                           { return nil }
func (m *MockDB) DeleteAllTasks() error                                    { return nil }
func (m *MockDB) FindSettings(settingsId string) (models.Settings, error) {
	return models.Settings{}, nil
}
func (m *MockDB) SaveSettings(s models.Settings) error { return nil }

func (m *MockDB) SaveTag(tagId string) error {
	if m.tags == nil {
		m.tags = make(map[string]bool)
	}
	m.tags[tagId] = true
	return nil
}

func (m *MockDB) TaskTags(taskId string) ([]models.TaskTag, error) {
	if m.taskTags == nil {
		return nil, nil
	}
	return m.taskTags[taskId], nil
}

func (m *MockDB) Tags() ([]models.TaskTag, error) {
	var tags []models.TaskTag
	for tag := range m.tags {
		tags = append(tags, models.TaskTag(tag))
	}
	return tags, nil
}

func (m *MockDB) AddTagToTask(taskId, tagId string) error {
	if m.taskTags == nil {
		m.taskTags = make(map[string][]models.TaskTag)
	}
	m.taskTags[taskId] = append(m.taskTags[taskId], models.TaskTag(tagId))
	return nil
}

func (m *MockDB) DeleteTagFromTask(taskId, tagId string) error {
	if m.taskTags == nil {
		return db.ErrNotFound
	}

	tags, exists := m.taskTags[taskId]
	if !exists {
		return db.ErrNotFound
	}

	for i, tag := range tags {
		if string(tag) == tagId {
			m.taskTags[taskId] = slices.Delete(tags, i, i+1)
			return nil
		}
	}

	return db.ErrNotFound
}

func (m *MockDB) FindTask(taskId string) (models.Task, error) {
	if task, exists := m.tasks[taskId]; exists {
		return task, nil
	}
	return models.Task{}, db.ErrNotFound
}

func (m *MockDB) SaveTask(task models.Task) error {
	if m.tasks == nil {
		m.tasks = make(map[string]models.Task)
	}
	m.tasks[task.Id] = task
	return nil
}

func (m *MockDB) MigrationExists(id string) bool {
	return m.migrations[id]
}

func (m *MockDB) RecordMigration(id string) {
	if m.migrations == nil {
		m.migrations = make(map[string]bool)
	}
	m.migrations[id] = true
}

func (m *MockDB) TasksTags(taskIds []string) (map[string][]models.TaskTag, error) {
	result := make(map[string][]models.TaskTag)
	if m.taskTags == nil {
		return result, nil
	}

	for _, taskId := range taskIds {
		if tags, exists := m.taskTags[taskId]; exists {
			result[taskId] = tags
		}
	}
	return result, nil
}

func setupTestDB() *MockDB {
	mockDB := &MockDB{
		tasks:      make(map[string]models.Task),
		migrations: make(map[string]bool),
		tags:       make(map[string]bool),
		taskTags:   make(map[string][]models.TaskTag),
	}
	db.SetDB(mockDB)
	return mockDB
}

func TestSaveNewTask_DoNotTouchTitleIfItDefined(t *testing.T) {
	mockDB := setupTestDB()

	task := models.Task{
		Title:     "ExistingTitle",
		Content:   "LongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitle",
		Priority:  models.PriorityHigh,
		Impact:    models.ImpactHigh,
		Cost:      models.CostL,
		Completed: models.NOT_COMPLETED,
	}

	err := SaveNewTask(task, nil)
	if err != nil {
		t.Errorf("SaveNewTask failed: %v", err)
	}

	actual, _ := mockDB.Tasks()
	expectedTitle := "ExistingTitle"
	if actual[0].Title != expectedTitle {
		t.Errorf("error; actual=%v; expected=%v", actual[0].Title, expectedTitle)
	}
}

func TestSaveNewTask_TitleCorrectlyUpdated(t *testing.T) {
	mockDB := setupTestDB()

	task := models.Task{
		Title:     "",
		Content:   "LongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitle",
		Priority:  models.PriorityHigh,
		Impact:    models.ImpactHigh,
		Cost:      models.CostL,
		Completed: models.NOT_COMPLETED,
	}

	err := SaveNewTask(task, nil)
	if err != nil {
		t.Errorf("SaveNewTask failed: %v", err)
	}

	actual, _ := mockDB.Tasks()
	expectedTitle := "LongTitleLongTitleLongTitleLongTitleLongTitleLongTitleLongTitleL"
	if actual[0].Title != expectedTitle {
		t.Errorf("error; actual=%v; expected=%v", actual[0].Title, expectedTitle)
	}
}

func TestSaveNewTask_TitleDoesNotHaveLineBreak(t *testing.T) {
	mockDB := setupTestDB()

	task := models.Task{
		Title:     "",
		Content:   "LongTitleLongTitleLong\nAfterBreakAfterBreakAfterBreakAfterBreakAfterBreakAfterBreakAfterBreakAfterBreakAfterBreak",
		Priority:  models.PriorityHigh,
		Impact:    models.ImpactHigh,
		Cost:      models.CostL,
		Completed: models.NOT_COMPLETED,
	}

	err := SaveNewTask(task, nil)
	if err != nil {
		t.Errorf("SaveNewTask failed: %v", err)
	}

	actual, _ := mockDB.Tasks()
	expectedTitle := "LongTitleLongTitleLong"
	if actual[0].Title != expectedTitle {
		t.Errorf("error; actual=%v; expected=%v", actual[0].Title, expectedTitle)
	}
}

func TestUpdateTask(t *testing.T) {
	mockDB := setupTestDB()

	taskId := "Id"
	originalTask := models.Task{
		Id:       taskId,
		Title:    "Title",
		Content:  "Content",
		Priority: models.PriorityLow,
	}
	mockDB.SaveTask(originalTask)

	updatedTask := originalTask
	updatedTask.Title = "NewTitle"
	updatedTask.Content = "NewContent"
	updatedTask.Priority = models.PriorityLow

	err := UpdateTask(updatedTask, nil)
	if err != nil {
		t.Errorf("UpdateTask failed: %v", err)
	}

	saved, err := mockDB.FindTask(taskId)
	if err != nil {
		t.Errorf("Failed to retrieve updated task: %v", err)
	}
	if saved.Title != updatedTask.Title ||
		saved.Content != updatedTask.Content ||
		saved.Priority != updatedTask.Priority {
		t.Error("Task was not properly updated")
	}
}

func TestSaveTask(t *testing.T) {
	mockDB := setupTestDB()

	task := models.Task{
		Id:       "test-id",
		Title:    "Test Task",
		Content:  "Test Content",
		Priority: models.PriorityHigh,
		Impact:   models.ImpactHigh,
		Cost:     models.CostL,
	}

	err := SaveTask(task)
	if err != nil {
		t.Errorf("SaveTask failed: %v", err)
	}

	saved, err := mockDB.FindTask("test-id")
	if err != nil {
		t.Errorf("Failed to retrieve saved task: %v", err)
	}

	fmt.Printf("saved=%v", saved)

	if saved.Id != task.Id ||
		saved.Title != task.Title ||
		saved.Content != task.Content ||
		saved.Priority != task.Priority ||
		saved.Impact != task.Impact ||
		saved.Cost != task.Cost ||
		saved.Value == 0.0 {
		t.Error("Task was not saved correctly")
	}
}

func TestUpdateTaskTags(t *testing.T) {
	mockDB := setupTestDB()
	taskId := "test-task"

	// Initial tags
	initialTags := []models.TaskTag{"tag1", "tag2"}
	for _, tag := range initialTags {
		mockDB.SaveTag(string(tag))
		mockDB.AddTagToTask(taskId, string(tag))
	}

	// Test cases
	tests := []struct {
		name        string
		changedTags []models.TaskTag
		wantErr     bool
	}{
		{
			name:        "Add new tag",
			changedTags: []models.TaskTag{"tag1", "tag2", "tag3"},
			wantErr:     false,
		},
		{
			name:        "Remove existing tag",
			changedTags: []models.TaskTag{"tag1"},
			wantErr:     false,
		},
		{
			name:        "No changes",
			changedTags: []models.TaskTag{"tag1", "tag2"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateTaskTags(taskId, tt.changedTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateTaskTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify tags were updated correctly
			resultTags, err := TaskTags(taskId)
			if err != nil {
				t.Errorf("Failed to get task tags: %v", err)
				return
			}

			if len(resultTags) != len(tt.changedTags) {
				t.Errorf("Expected %d tags, got %d", len(tt.changedTags), len(resultTags))
			}

			// Check if all expected tags are present
			for _, expectedTag := range tt.changedTags {
				found := slices.Contains(resultTags, expectedTag)
				if !found {
					t.Errorf("Expected tag %s not found in result", expectedTag)
				}
			}
		})
	}
}
