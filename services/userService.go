package services

import (
	"errors"
	"fmt"

	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

const SETTINGS_ID = "UserSettings"

func SetCompletedFilter(val bool) error {
	s, err := FindUserSettings()
	if err != nil {
		return err
	}
	s.FilterCompleted = val
	err = UpdateUserSettings(s)
	return err
}

func FindUserSettings() (models.Settings, error) {
	var s models.Settings
	var err error

	s, err = db.FindSettings(SETTINGS_ID)
	if errors.Is(err, db.ErrNotFound) {
		err = UpdateUserSettings(models.Settings{
			Id:              SETTINGS_ID,
			FilterCompleted: true,
		})
	}
	if err != nil {
		err = fmt.Errorf("failed to retrieve settings: %w", err)
	}

	return s, err
}

func UpdateUserSettings(s models.Settings) error {
	return db.SaveSettings(s)
}

func ToggleSorting(s models.Settings, newColumn models.SortColumn, actDir models.SortDirection) error {
	if s.ActiveSortColumn == newColumn {
		s.ActiveSortDirection = actDir.Flip()
	} else {
		s.ActiveSortDirection = models.Desc // default
	}
	s.ActiveSortColumn = newColumn
	return UpdateUserSettings(s)
}
