package services

import (
	"testing"

	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

// MockDB implements db.Db interface for testing
type MockDB struct {
	tasks map[string]models.Task
}

func (m *MockDB) Init()  {}
func (m *MockDB) Close() {}
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

func setupTestDB() *MockDB {
	mockDB := &MockDB{
		tasks: make(map[string]models.Task),
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

	err := SaveNewTask(task)
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

	err := SaveNewTask(task)
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

	err := SaveNewTask(task)
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

	err := UpdateTask(updatedTask)
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
