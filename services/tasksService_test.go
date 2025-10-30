package services

import (
	"fmt"
	"strings"
	"testing"

	"slices"

	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

type MockDB struct {
	db.NoOpDB
	tasks      map[string]models.Task
	migrations map[string]bool
	tags       map[string]bool
	taskTags   map[string][]models.TaskTag
}

func (m *MockDB) Tasks() ([]models.Task, error) {
	var result []models.Task
	for _, value := range m.tasks {
		result = append(result, value)
	}
	return result, nil
}

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

type reducePriorityMockDB struct {
	db.NoOpDB
	expectedTasks []models.Task
	savedTasks    []models.Task
}

func (m *reducePriorityMockDB) FindTasks(query models.TasksQuery) ([]models.Task, error) {
	return m.expectedTasks, nil
}

func (m *reducePriorityMockDB) SaveTask(task models.Task) error {
	m.savedTasks = append(m.savedTasks, task)
	return nil
}

func TestReducePriorityForVisibleTasks(t *testing.T) {
	tests := []struct {
		name          string
		tasks         []models.Task
		query         models.TasksQuery
		expectedSaves map[string]models.TaskPriority // map of task ID to expected priority after save
	}{
		{
			name: "reduce priorities for all visible tasks",
			tasks: []models.Task{
				{Id: "1", Priority: models.PriorityUrgent},
				{Id: "2", Priority: models.PriorityHigh},
				{Id: "3", Priority: models.PriorityMedium},
				{Id: "4", Priority: models.PriorityLow},
			},
			query: models.TasksQuery{},
			expectedSaves: map[string]models.TaskPriority{
				"1": models.PriorityHigh,
				"2": models.PriorityMedium,
				"3": models.PriorityLow,
			},
		},
		{
			name:          "empty task list",
			tasks:         []models.Task{},
			query:         models.TasksQuery{},
			expectedSaves: map[string]models.TaskPriority{},
		},
		{
			name: "all tasks already at minimum priority",
			tasks: []models.Task{
				{Id: "1", Priority: models.PriorityLow},
				{Id: "2", Priority: models.PriorityLow},
			},
			query:         models.TasksQuery{},
			expectedSaves: map[string]models.TaskPriority{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &reducePriorityMockDB{
				expectedTasks: tt.tasks,
			}
			db.SetDB(mockDB)

			err := ReducePriorityForVisibleTasks(tt.query)
			if err != nil {
				t.Fatalf("ReducePriorityForVisibleTasks failed: %v", err)
			}

			// Verify that only tasks that needed priority reduction were saved
			if len(mockDB.savedTasks) != len(tt.expectedSaves) {
				t.Errorf("Expected %d saves, got %d", len(tt.expectedSaves), len(mockDB.savedTasks))
			}

			// Create a map of saved tasks by ID for easier lookup
			savedTasksMap := make(map[string]models.Task)
			for _, task := range mockDB.savedTasks {
				savedTasksMap[task.Id] = task
			}

			// Verify each expected save
			for taskID, expectedPriority := range tt.expectedSaves {
				savedTask, exists := savedTasksMap[taskID]
				if !exists {
					t.Errorf("Expected task %s to be saved with priority %v, but it wasn't saved",
						taskID, expectedPriority)
					continue
				}

				if savedTask.Priority != expectedPriority {
					t.Errorf("Task %s: expected saved priority %v, got %v",
						taskID, expectedPriority, savedTask.Priority)
				}
			}

			// Verify no unexpected saves occurred
			for _, savedTask := range mockDB.savedTasks {
				expectedPriority, shouldBeSaved := tt.expectedSaves[savedTask.Id]
				if !shouldBeSaved {
					t.Errorf("Task %s was saved but shouldn't have been", savedTask.Id)
					continue
				}

				if savedTask.Priority != expectedPriority {
					t.Errorf("Task %s: wrong priority saved. Expected %v, got %v",
						savedTask.Id, expectedPriority, savedTask.Priority)
				}
			}
		})
	}
}

// ai: Happy Path Tests for CloneTask
func Test_CloneTask_Success(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:       "original-id",
		Title:    "Original Task",
		Content:  "Original content",
		Priority: models.PriorityHigh,
		Impact:   models.ImpactModerate,
		Cost:     models.CostM,
		Fun:      models.FunM,
		Wip:      true,
		Planned:  false,
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if clonedTask.Id == originalTask.Id {
		t.Error("Cloned task should have different ID")
	}
	if clonedTask.Title != "Copy of Original Task" {
		t.Errorf("Expected title 'Copy of Original Task', got '%s'", clonedTask.Title)
	}
	if clonedTask.Content != originalTask.Content {
		t.Error("Content should be copied")
	}
	if clonedTask.Priority != originalTask.Priority {
		t.Error("Priority should be copied")
	}
	if clonedTask.Impact != originalTask.Impact {
		t.Error("Impact should be copied")
	}
	if clonedTask.Cost != originalTask.Cost {
		t.Error("Cost should be copied")
	}
	if clonedTask.Fun != originalTask.Fun {
		t.Error("Fun should be copied")
	}
	if clonedTask.Wip != originalTask.Wip {
		t.Error("Wip should be copied")
	}
	if clonedTask.Planned != originalTask.Planned {
		t.Error("Planned should be copied")
	}
}

func Test_CloneTask_WithTags(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:    "original-id",
		Title: "Task with tags",
	}
	mockDB.SaveTask(originalTask)

	// ai: Add tags to original task
	tags := []models.TaskTag{"urgent", "work", "project"}
	for _, tag := range tags {
		mockDB.SaveTag(string(tag))
		mockDB.AddTagToTask("original-id", string(tag))
	}

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if len(clonedTask.Tags) != len(tags) {
		t.Errorf("Expected %d tags, got %d", len(tags), len(clonedTask.Tags))
	}

	for _, expectedTag := range tags {
		found := false
		for _, actualTag := range clonedTask.Tags {
			if actualTag == expectedTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected tag '%s' not found in cloned task", expectedTag)
		}
	}
}

func Test_CloneTask_WithoutTags(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:    "original-id",
		Title: "Task without tags",
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if len(clonedTask.Tags) != 0 {
		t.Errorf("Expected 0 tags, got %d", len(clonedTask.Tags))
	}
}

func Test_CloneTask_CompletedTask(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:        "original-id",
		Title:     "Completed Task",
		Completed: models.NOT_COMPLETED.Add(24 * 60 * 60 * 1000000000), // 1 day ago
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if clonedTask.Completed != models.NOT_COMPLETED {
		t.Error("Cloned task should have completion status reset")
	}
}

// ai: Property Validation Tests for CloneTask
func Test_CloneTask_TitlePrefix(t *testing.T) {
	mockDB := setupTestDB()

	testCases := []struct {
		name          string
		originalTitle string
		expectedTitle string
	}{
		{"Simple title", "My Task", "Copy of My Task"},
		{"Empty title", "", "Copy of "},
		{"Title with spaces", "  Task with spaces  ", "Copy of   Task with spaces  "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			taskId := "task-" + tc.name
			originalTask := models.Task{
				Id:    taskId,
				Title: tc.originalTitle,
			}
			mockDB.SaveTask(originalTask)

			clonedTask, err := CloneTask(taskId)
			if err != nil {
				t.Fatalf("CloneTask failed: %v", err)
			}

			if clonedTask.Title != tc.expectedTitle {
				t.Errorf("Expected title '%s', got '%s'", tc.expectedTitle, clonedTask.Title)
			}
		})
	}
}

func Test_CloneTask_NewIdGenerated(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:    "original-id",
		Title: "Test Task",
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if clonedTask.Id == originalTask.Id {
		t.Error("Cloned task should have different ID than original")
	}
	if clonedTask.Id == "" {
		t.Error("Cloned task should have non-empty ID")
	}
}

func Test_CloneTask_NewTimestamps(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:      "original-id",
		Title:   "Test Task",
		Created: models.NOT_COMPLETED.Add(-24 * 60 * 60 * 1000000000), // 1 day ago
		Updated: models.NOT_COMPLETED.Add(-12 * 60 * 60 * 1000000000), // 12 hours ago
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if clonedTask.Created.Equal(originalTask.Created) {
		t.Error("Cloned task should have different Created timestamp")
	}
	// ai: Note: AsNewTask() only sets Created timestamp, not Updated
}

func Test_CloneTask_CompletionReset(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:        "original-id",
		Title:     "Completed Task",
		Completed: models.NOT_COMPLETED.Add(24 * 60 * 60 * 1000000000), // 1 day ago
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if clonedTask.Completed != models.NOT_COMPLETED {
		t.Errorf("Expected completion to be reset to NOT_COMPLETED, got %v", clonedTask.Completed)
	}
}

func Test_CloneTask_OtherPropertiesCopied(t *testing.T) {
	mockDB := setupTestDB()

	originalTask := models.Task{
		Id:       "original-id",
		Title:    "Test Task",
		Content:  "Test content with details",
		Priority: models.PriorityUrgent,
		Impact:   models.ImpactHigh,
		Cost:     models.CostXL,
		Fun:      models.FunL,
		Wip:      true,
		Planned:  true,
	}
	mockDB.SaveTask(originalTask)

	clonedTask, err := CloneTask("original-id")
	if err != nil {
		t.Fatalf("CloneTask failed: %v", err)
	}

	if clonedTask.Content != originalTask.Content {
		t.Errorf("Expected content '%s', got '%s'", originalTask.Content, clonedTask.Content)
	}
	if clonedTask.Priority != originalTask.Priority {
		t.Errorf("Expected priority %v, got %v", originalTask.Priority, clonedTask.Priority)
	}
	if clonedTask.Impact != originalTask.Impact {
		t.Errorf("Expected impact %v, got %v", originalTask.Impact, clonedTask.Impact)
	}
	if clonedTask.Cost != originalTask.Cost {
		t.Errorf("Expected cost %v, got %v", originalTask.Cost, clonedTask.Cost)
	}
	if clonedTask.Fun != originalTask.Fun {
		t.Errorf("Expected fun %v, got %v", originalTask.Fun, clonedTask.Fun)
	}
	if clonedTask.Wip != originalTask.Wip {
		t.Errorf("Expected wip %v, got %v", originalTask.Wip, clonedTask.Wip)
	}
	if clonedTask.Planned != originalTask.Planned {
		t.Errorf("Expected planned %v, got %v", originalTask.Planned, clonedTask.Planned)
	}
}

// ai: Error Handling Tests for CloneTask
func Test_CloneTask_TaskNotFound(t *testing.T) {
	setupTestDB()

	_, err := CloneTask("non-existent-id")
	if err == nil {
		t.Error("Expected error when cloning non-existent task")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to find original task") {
		t.Errorf("Expected 'failed to find original task' error, got: %v", err)
	}
}

// ai: Mock for error testing
type errorMockDB struct {
	MockDB
	findTaskError error
	taskTagsError error
	saveTaskError error
	addTagError   error
}

func (m *errorMockDB) FindTask(taskId string) (models.Task, error) {
	if m.findTaskError != nil {
		return models.Task{}, m.findTaskError
	}
	return m.MockDB.FindTask(taskId)
}

func (m *errorMockDB) TaskTags(taskId string) ([]models.TaskTag, error) {
	if m.taskTagsError != nil {
		return nil, m.taskTagsError
	}
	return m.MockDB.TaskTags(taskId)
}

func (m *errorMockDB) SaveTask(task models.Task) error {
	if m.saveTaskError != nil {
		return m.saveTaskError
	}
	return m.MockDB.SaveTask(task)
}

func (m *errorMockDB) AddTagToTask(taskId, tagId string) error {
	if m.addTagError != nil {
		return m.addTagError
	}
	return m.MockDB.AddTagToTask(taskId, tagId)
}

func Test_CloneTask_DatabaseError(t *testing.T) {
	mockDB := &errorMockDB{
		MockDB: MockDB{
			tasks: map[string]models.Task{
				"test-id": {Id: "test-id", Title: "Test"},
			},
		},
		taskTagsError: fmt.Errorf("database connection error"),
	}
	db.SetDB(mockDB)

	_, err := CloneTask("test-id")
	if err == nil {
		t.Error("Expected error when database operation fails")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to get original task tags") {
		t.Errorf("Expected 'failed to get original task tags' error, got: %v", err)
	}
}

func Test_CloneTask_SaveTaskFails(t *testing.T) {
	mockDB := &errorMockDB{
		MockDB: MockDB{
			tasks: map[string]models.Task{
				"test-id": {Id: "test-id", Title: "Test"},
			},
			taskTags: map[string][]models.TaskTag{
				"test-id": {},
			},
		},
		saveTaskError: fmt.Errorf("save operation failed"),
	}
	db.SetDB(mockDB)

	_, err := CloneTask("test-id")
	if err == nil {
		t.Error("Expected error when SaveTask fails")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to save cloned task") {
		t.Errorf("Expected 'failed to save cloned task' error, got: %v", err)
	}
}

func Test_CloneTask_AddTagFails(t *testing.T) {
	mockDB := &errorMockDB{
		MockDB: MockDB{
			tasks: map[string]models.Task{
				"test-id": {Id: "test-id", Title: "Test"},
			},
			taskTags: map[string][]models.TaskTag{
				"test-id": {"tag1", "tag2"},
			},
		},
		addTagError: fmt.Errorf("add tag operation failed"),
	}
	db.SetDB(mockDB)

	_, err := CloneTask("test-id")
	if err == nil {
		t.Error("Expected error when AddTagToTask fails")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to add tag to cloned task") {
		t.Errorf("Expected 'failed to add tag to cloned task' error, got: %v", err)
	}
}

func Test_CloneTask_GetTagsFails(t *testing.T) {
	mockDB := &errorMockDB{
		MockDB: MockDB{
			tasks: map[string]models.Task{
				"test-id": {Id: "test-id", Title: "Test"},
			},
		},
		taskTagsError: fmt.Errorf("get tags operation failed"),
	}
	db.SetDB(mockDB)

	_, err := CloneTask("test-id")
	if err == nil {
		t.Error("Expected error when TaskTags fails")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to get original task tags") {
		t.Errorf("Expected 'failed to get original task tags' error, got: %v", err)
	}
}
