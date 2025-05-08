package services

import (
	"testing"
	"time"

	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

type userSettingsTestDB struct {
	db.NoOpDB
	settings      models.Settings
	saveCallCount int
	findCallCount int
}

func (m *userSettingsTestDB) FindSettings(settingsId string) (models.Settings, error) {
	m.findCallCount++
	return m.settings, nil
}

func (m *userSettingsTestDB) SaveSettings(s models.Settings) error {
	m.saveCallCount++
	m.settings = s
	return nil
}

func setupUserSettingsTestDB() *userSettingsTestDB {
	mockDB := &userSettingsTestDB{
		settings: models.Settings{
			Id: SETTINGS_ID,
			TasksQuery: models.TasksQuery{
				FilterCompleted: true,
			},
		},
	}
	db.SetDB(mockDB)
	return mockDB
}

func Test_SetCompletedFilter_Success(t *testing.T) {
	mockDB := setupUserSettingsTestDB()

	err := SetCompletedFilter(false)
	if err != nil {
		t.Errorf("SetCompletedFilter failed: %v", err)
	}

	if mockDB.settings.TasksQuery.FilterCompleted != false {
		t.Error("FilterCompleted was not updated correctly")
	}

	if mockDB.saveCallCount != 1 {
		t.Errorf("Expected 1 save call, got %d", mockDB.saveCallCount)
	}

	if mockDB.findCallCount != 1 {
		t.Errorf("Expected 1 find call, got %d", mockDB.findCallCount)
	}
}

func Test_RemoveTagFromSettings_Success(t *testing.T) {
	mockDB := setupUserSettingsTestDB()
	mockDB.settings.TasksQuery.Tags = []models.TaskTag{"tag1", "tag2", "tag3"}

	err := RemoveTagFromSettings("tag2")
	if err != nil {
		t.Errorf("RemoveTagFromSettings failed: %v", err)
	}

	if len(mockDB.settings.TasksQuery.Tags) != 2 {
		t.Errorf("Expected 2 tags after removal, got %d", len(mockDB.settings.TasksQuery.Tags))
	}

	for _, tag := range mockDB.settings.TasksQuery.Tags {
		if tag == "tag2" {
			t.Error("Tag was not removed correctly")
		}
	}
}

func Test_ToggleSorting_NewColumn(t *testing.T) {
	mockDB := setupUserSettingsTestDB()
	settings := mockDB.settings

	err := ToggleSorting(settings, models.Priority, models.Asc)
	if err != nil {
		t.Errorf("ToggleSorting failed: %v", err)
	}

	if mockDB.settings.TasksQuery.SortColumn != models.Priority {
		t.Error("SortColumn was not updated correctly")
	}

	if mockDB.settings.TasksQuery.SortDirection != models.Desc {
		t.Error("SortDirection was not set to default Desc")
	}
}

func Test_ToggleSorting_SameColumn(t *testing.T) {
	mockDB := setupUserSettingsTestDB()
	settings := mockDB.settings
	settings.TasksQuery.SortColumn = models.Priority
	settings.TasksQuery.SortDirection = models.Desc

	err := ToggleSorting(settings, models.Priority, models.Asc)
	if err != nil {
		t.Errorf("ToggleSorting failed: %v", err)
	}

	if mockDB.settings.TasksQuery.SortColumn != models.Priority {
		t.Error("SortColumn should not change")
	}

	if mockDB.settings.TasksQuery.SortDirection != models.Desc {
		t.Error("SortDirection was not flipped correctly")
	}
}

func Test_ApplyPreparedQuery_CompletedToday(t *testing.T) {
	mockDB := setupUserSettingsTestDB()

	err := ApplyPreparedQuery(consts.PREPARED_QUERY_COMPLETED_TODAY)
	if err != nil {
		t.Errorf("ApplyPreparedQuery failed: %v", err)
	}

	q := mockDB.settings.TasksQuery
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	nextMidnight := midnight.AddDate(0, 0, 1)

	if !q.CompletedFrom.Equal(midnight) {
		t.Errorf("CompletedFrom not set correctly, got %v, want %v", q.CompletedFrom, midnight)
	}

	if !q.CompletedTo.Equal(nextMidnight) {
		t.Errorf("CompletedTo not set correctly, got %v, want %v", q.CompletedTo, nextMidnight)
	}

	if q.FilterCompleted || !q.FilterIncompleted {
		t.Error("Filter flags not set correctly")
	}

	if q.SortColumn != models.Completed {
		t.Error("Sort column not set to Completed")
	}
}

func Test_ApplyPreparedQuery_CompletedYesterday(t *testing.T) {
	mockDB := setupUserSettingsTestDB()

	err := ApplyPreparedQuery(consts.PREPARED_QUERY_COMPLETED_YESTERDAY)
	if err != nil {
		t.Errorf("ApplyPreparedQuery failed: %v", err)
	}

	q := mockDB.settings.TasksQuery
	now := time.Now()
	midnightYesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	midnightToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if !q.CompletedFrom.Equal(midnightYesterday) {
		t.Errorf("CompletedFrom not set correctly, got %v, want %v", q.CompletedFrom, midnightYesterday)
	}

	if !q.CompletedTo.Equal(midnightToday) {
		t.Errorf("CompletedTo not set correctly, got %v, want %v", q.CompletedTo, midnightToday)
	}

	if q.FilterCompleted || !q.FilterIncompleted {
		t.Error("Filter flags not set correctly")
	}

	if q.SortColumn != models.Completed {
		t.Error("Sort column not set to Completed")
	}
}

func Test_ApplyPreparedQuery_Reset(t *testing.T) {
	mockDB := setupUserSettingsTestDB()
	mockDB.settings.TasksQuery.Tags = []models.TaskTag{"tag1", "tag2"}
	mockDB.settings.TasksQuery.FilterCompleted = false
	mockDB.settings.TasksQuery.FilterIncompleted = true
	mockDB.settings.TasksQuery.FilterWip = true
	mockDB.settings.TasksQuery.SortColumn = models.Priority

	err := ApplyPreparedQuery(consts.PREPARED_QUERY_RESET)
	if err != nil {
		t.Errorf("ApplyPreparedQuery failed: %v", err)
	}

	q := mockDB.settings.TasksQuery
	if len(q.Tags) != 0 {
		t.Error("Tags were not reset")
	}
	if !q.FilterCompleted {
		t.Error("FilterCompleted not reset to true")
	}
	if q.FilterIncompleted {
		t.Error("FilterIncompleted not reset to false")
	}
	if q.FilterWip {
		t.Error("FilterWip not reset to false")
	}
	if q.SortColumn != models.Priority {
		t.Error("SortColumn not set to Priority")
	}
	if q.SortDirection != models.Desc {
		t.Error("SortDirection not set to Desc")
	}
}

func Test_ApplyPreparedQuery_CompletedThisWeek(t *testing.T) {
	mockDB := setupUserSettingsTestDB()

	err := ApplyPreparedQuery(consts.PREPARED_QUERY_COMPLETED_THIS_WEEK)
	if err != nil {
		t.Errorf("ApplyPreparedQuery failed: %v", err)
	}

	q := mockDB.settings.TasksQuery
	mondayMidnight := thisMonday()
	nextMonday := mondayMidnight.AddDate(0, 0, 7)

	if !q.CompletedFrom.Equal(mondayMidnight) {
		t.Errorf("CompletedFrom not set correctly, got %v, want %v", q.CompletedFrom, mondayMidnight)
	}

	if !q.CompletedTo.Equal(nextMonday) {
		t.Errorf("CompletedTo not set correctly, got %v, want %v", q.CompletedTo, nextMonday)
	}

	if q.FilterCompleted || !q.FilterIncompleted {
		t.Error("Filter flags not set correctly")
	}

	if q.SortColumn != models.Completed {
		t.Error("Sort column not set to Completed")
	}
}

func Test_ApplyPreparedQuery_CompletedLastWeek(t *testing.T) {
	mockDB := setupUserSettingsTestDB()

	err := ApplyPreparedQuery(consts.PREPARED_QUERY_COMPLETED_LAST_WEEK)
	if err != nil {
		t.Errorf("ApplyPreparedQuery failed: %v", err)
	}

	q := mockDB.settings.TasksQuery
	previousMonday := thisMonday().AddDate(0, 0, -7)

	if !q.CompletedFrom.Equal(previousMonday) {
		t.Errorf("CompletedFrom not set correctly, got %v, want %v", q.CompletedFrom, previousMonday)
	}

	if q.FilterCompleted || !q.FilterIncompleted {
		t.Error("Filter flags not set correctly")
	}

	if q.SortColumn != models.Completed {
		t.Error("Sort column not set to Completed")
	}
}

func Test_ApplyPreparedQuery_CompletedLastTwoWeeks(t *testing.T) {
	mockDB := setupUserSettingsTestDB()

	err := ApplyPreparedQuery(consts.PREPARED_QUERY_COMPLETED_LAST_TWO_WEEKS)
	if err != nil {
		t.Errorf("ApplyPreparedQuery failed: %v", err)
	}

	q := mockDB.settings.TasksQuery
	now := time.Now()
	twoWeeksAgo := time.Date(now.Year(), now.Month(), now.Day()-14, 0, 0, 0, 0, now.Location())

	if !q.CompletedFrom.Equal(twoWeeksAgo) {
		t.Errorf("CompletedFrom not set correctly, got %v, want %v", q.CompletedFrom, twoWeeksAgo)
	}

	if q.FilterCompleted || !q.FilterIncompleted {
		t.Error("Filter flags not set correctly")
	}

	if q.SortColumn != models.Completed {
		t.Error("Sort column not set to Completed")
	}
}

func Test_ApplyPreparedQuery_InvalidQueryName(t *testing.T) {
	mockDB := setupUserSettingsTestDB()
	originalSettings := mockDB.settings

	err := ApplyPreparedQuery("invalid-query-name")
	if err != nil {
		t.Errorf("ApplyPreparedQuery should not return error for invalid query name, got: %v", err)
	}

	// Settings should remain unchanged for invalid query name
	if mockDB.settings.TasksQuery.FilterCompleted != originalSettings.TasksQuery.FilterCompleted {
		t.Error("Settings were modified when they should not have been")
	}
}

func Test_UpdateUserSettings_Success(t *testing.T) {
	mockDB := setupUserSettingsTestDB()
	newSettings := models.Settings{
		Id: SETTINGS_ID,
		TasksQuery: models.TasksQuery{
			FilterCompleted: false,
			FilterWip:       true,
			Tags:            []models.TaskTag{"tag1", "tag2"},
		},
	}

	err := UpdateUserSettings(newSettings)
	if err != nil {
		t.Errorf("UpdateUserSettings failed: %v", err)
	}

	if mockDB.saveCallCount != 1 {
		t.Errorf("Expected 1 save call, got %d", mockDB.saveCallCount)
	}

	// Verify all fields were updated
	if mockDB.settings.TasksQuery.FilterCompleted != false {
		t.Error("FilterCompleted was not updated")
	}
	if mockDB.settings.TasksQuery.FilterWip != true {
		t.Error("FilterWip was not updated")
	}
	if len(mockDB.settings.TasksQuery.Tags) != 2 {
		t.Error("Tags were not updated correctly")
	}
}
