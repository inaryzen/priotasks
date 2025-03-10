package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

const SETTINGS_ID = "UserSettings"

func SetCompletedFilter(val bool) error {
	s, err := FindUserSettings()
	if err != nil {
		return err
	}
	s.TasksQuery.FilterCompleted = val
	err = UpdateUserSettings(s)
	return err
}

func FindUserSettings() (models.Settings, error) {
	var s models.Settings
	var err error

	s, err = db.DB().FindSettings(SETTINGS_ID)
	if errors.Is(err, db.ErrNotFound) {
		err = UpdateUserSettings(models.Settings{
			Id: SETTINGS_ID,
			TasksQuery: models.TasksQuery{
				FilterCompleted: true,
			},
		})
	}
	if err != nil {
		err = fmt.Errorf("failed to retrieve settings: %w", err)
	}

	return s, err
}

func UpdateUserSettings(s models.Settings) error {
	return db.DB().SaveSettings(s)
}

func ToggleSorting(s models.Settings, newColumn models.SortColumn, actDir models.SortDirection) error {
	if s.TasksQuery.SortColumn == newColumn {
		s.TasksQuery.SortDirection = actDir.Flip()
	} else {
		s.TasksQuery.SortDirection = models.Desc // default
	}
	s.TasksQuery.SortColumn = newColumn
	return UpdateUserSettings(s)
}

func ApplyPreparedQuery(preparedQueryName string) error {
	s, err := FindUserSettings()
	if err != nil {
		return err
	}
	q := s.TasksQuery
	q = q.Reset()

	common.Debug("ApplyPreparedQuery: %v", preparedQueryName)

	switch preparedQueryName {
	case consts.PREPARED_QUERY_COMPLETED_YESTERDAY:
		now := time.Now()
		midnightYesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		q.CompletedFrom = midnightYesterday
		q.FilterIncompleted = true
		q.FilterCompleted = false
		q.SortColumn = models.Completed
	case consts.PREPARED_QUERY_COMPLETED_TODAY:
		now := time.Now()
		midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		q.CompletedFrom = midnight
		q.FilterIncompleted = true
		q.FilterCompleted = false
		q.SortColumn = models.Completed
	case consts.PREPARED_QUERY_COMPLETED_THIS_WEEK:
		mondayMidnight := thisMonday()
		q.CompletedFrom = mondayMidnight
		q.FilterIncompleted = true
		q.FilterCompleted = false
		q.SortColumn = models.Completed
	case consts.PREPARED_QUERY_COMPLETED_LAST_TWO_WEEKS:
		now := time.Now()
		twoWeeksAgo := time.Date(now.Year(), now.Month(), now.Day()-14, 0, 0, 0, 0, now.Location())
		q.CompletedFrom = twoWeeksAgo
		q.FilterIncompleted = true
		q.FilterCompleted = false
		q.SortColumn = models.Completed
	case consts.PREPARED_QUERY_COMPLETED_LAST_WEEK:
		previousMonday := thisMonday().AddDate(0, 0, -7)
		q.CompletedFrom = previousMonday
		q.FilterIncompleted = true
		q.FilterCompleted = false
		q.SortColumn = models.Completed
	// case PREPARED_QUERY_RESET:
	// nop
	default:
		// nop
	}

	common.Debug("ApplyPreparedQuery: %v", q)
	s.TasksQuery = q
	UpdateUserSettings(s)

	return nil
}

func thisMonday() time.Time {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 { // time.Sunday == 0
		offset -= 7 // 1 - 0 - 7 == -6
	}
	return time.Date(now.Year(), now.Month(), now.Day()+offset, 0, 0, 0, 0, now.Location())
}
