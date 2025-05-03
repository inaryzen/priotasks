package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

type MockDB struct {
	settings models.Settings
}

func (m *MockDB) Init(p string)                                                   {}
func (m *MockDB) Close()                                                          {}
func (m *MockDB) Tasks() ([]models.Task, error)                                   { return nil, nil }
func (m *MockDB) FindTasks(query models.TasksQuery) ([]models.Task, error)        { return nil, nil }
func (m *MockDB) DeleteTask(taskId string) error                                  { return nil }
func (m *MockDB) DeleteAllTasks() error                                           { return nil }
func (m *MockDB) FindTask(taskId string) (models.Task, error)                     { return models.Task{}, nil }
func (m *MockDB) SaveTask(task models.Task) error                                 { return nil }
func (m *MockDB) SaveTag(tagId string) error                                      { return nil }
func (m *MockDB) AddTagToTask(taskId, tagId string) error                         { return nil }
func (m *MockDB) DeleteTagFromTask(taskId, tagId string) error                    { return nil }
func (m *MockDB) TaskTags(taskId string) ([]models.TaskTag, error)                { return nil, nil }
func (m *MockDB) Tags() ([]models.TaskTag, error)                                 { return nil, nil }
func (m *MockDB) TasksTags(taskIds []string) (map[string][]models.TaskTag, error) { return nil, nil }
func (m *MockDB) MigrationExists(id string) bool                                  { return false }
func (m *MockDB) RecordMigration(id string)                                       {}

func (m *MockDB) DeleteTag(tagId string) error             { return nil }
func (m *MockDB) DeleteTagFromAllTasks(tagId string) error { return nil }

func (m *MockDB) FindSettings(settingsId string) (models.Settings, error) {
	return m.settings, nil
}

func (m *MockDB) SaveSettings(s models.Settings) error {
	m.settings = s
	return nil
}

func setupTestHandler() *MockDB {
	mockDB := &MockDB{
		settings: models.Settings{
			Id: "UserSettings",
			TasksQuery: models.TasksQuery{
				FilterCompleted: false,
				CompletedFrom:   models.NOT_COMPLETED,
				CompletedTo:     models.NOT_COMPLETED,
				Tags:            []models.TaskTag{},
			},
		},
	}
	db.SetDB(mockDB)
	return mockDB
}

func TestPostFilterName_EmptyTag(t *testing.T) {
	mockDB := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/filter/"+consts.FILTER_TAGS, strings.NewReader(consts.FILTER_TAGS+"="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("name", consts.FILTER_TAGS)
	rr := httptest.NewRecorder()

	PostFilterName(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if len(mockDB.settings.TasksQuery.Tags) != 0 {
		t.Error("Tags should not have been added")
	}
}

func TestPostFilterName_ValidTag(t *testing.T) {
	mockDB := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/filter/"+consts.FILTER_TAGS, strings.NewReader(consts.FILTER_TAGS+"=test-tag"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("name", consts.FILTER_TAGS)
	rr := httptest.NewRecorder()

	PostFilterName(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if len(mockDB.settings.TasksQuery.Tags) != 1 {
		t.Error("Tag should have been added")
	}
	if len(mockDB.settings.TasksQuery.Tags) > 0 && mockDB.settings.TasksQuery.Tags[0] != "test-tag" {
		t.Errorf("Expected tag 'test-tag', got %s", mockDB.settings.TasksQuery.Tags[0])
	}
}

func TestPostFilterName_CompletedFilter(t *testing.T) {
	mockDB := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/filter/"+consts.FILTER_NAME_HIDE_COMPLETED,
		strings.NewReader(consts.FILTER_NAME_HIDE_COMPLETED+"=true"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("name", consts.FILTER_NAME_HIDE_COMPLETED)
	rr := httptest.NewRecorder()

	PostFilterName(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !mockDB.settings.TasksQuery.FilterCompleted {
		t.Error("FilterCompleted should be true")
	}
}

func TestPostFilterName_InvalidFilter(t *testing.T) {
	// mockDB := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/filter/invalid-filter", strings.NewReader("value=test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	PostFilterName(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}
