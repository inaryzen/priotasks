package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/models"
	_ "modernc.org/sqlite"
)

const (
	SETTINGS_COLUMNS = "id, filter_completed, filter_incompleted, active_sort_column, active_sort_direction, completed_from, completed_to, filter_wip, filter_non_wip, planned, non_planned, tags, search_text, enable_limit, limit_count"
)

func (d *DbSQLite) initSettings() {
	common.Debug("initSettings")
	_, err := d.instance.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id TEXT PRIMARY KEY,
			filter_completed BOOLEAN,
			filter_incompleted BOOLEAN,
			active_sort_column INTEGER,
			active_sort_direction INTEGER,
			completed_from TEXT,
			completed_to TEXT,
			filter_wip BOOLEAN,
			filter_non_wip BOOLEAN,
			planned BOOLEAN,
			non_planned BOOLEAN
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	d.addSettingsCompletedFrom()
	d.addSettingsCompletedTo()
	d.addSettingsFilterIncomplete()
	d.addSettingsFilterWipAndNonWip()
	d.addSettingsPlannedAndNonPlanned()
	d.settingsTableAddTagsColumn()
	d.settingsTableAddSearchTextColumn()
	d.settingsTableAddLimitColumns()
}

func (d *DbSQLite) settingsTableAddTagsColumn() {
	id := "settings_table_add_tags_column"
	if !d.MigrationExists(id) {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN tags TEXT DEFAULT ''")
		if err != nil {
			panic(err)
		} else {
			d.RecordMigration(id)
		}
	}
}

func (d *DbSQLite) settingsTableAddSearchTextColumn() {
	id := "settings_table_add_search_text_column"
	if !d.MigrationExists(id) {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN search_text TEXT DEFAULT ''")
		if err != nil {
			panic(err)
		} else {
			d.RecordMigration(id)
		}
	}
}

func (d *DbSQLite) settingsTableAddLimitColumns() {
	enableLimitId := "settings_table_add_enable_limit_column"
	if !d.MigrationExists(enableLimitId) {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN enable_limit BOOLEAN DEFAULT 1")
		if err != nil {
			panic(err)
		} else {
			d.RecordMigration(enableLimitId)
		}
	}

	limitCountId := "settings_table_add_limit_count_column"
	if !d.MigrationExists(limitCountId) {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN limit_count INTEGER DEFAULT 10")
		if err != nil {
			panic(err)
		} else {
			d.RecordMigration(limitCountId)
		}
	}
}

func (d *DbSQLite) addSettingsCompletedFrom() {
	if !d.columnExists("settings", "completed_from") {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN completed_from TEXT DEFAULT '" + models.NOT_COMPLETED.Format(consts.DEFAULT_DATE_FORMAT) + "'")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) addSettingsCompletedTo() {
	if !d.columnExists("settings", "completed_to") {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN completed_to TEXT DEFAULT '" + models.NOT_COMPLETED.Format(consts.DEFAULT_DATE_FORMAT) + "'")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) addSettingsFilterIncomplete() {
	if !d.columnExists("settings", "filter_incompleted") {
		_, err := d.instance.Exec("ALTER TABLE settings ADD COLUMN filter_incompleted BOOLEAN DEFAULT 0")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) addSettingsFilterWipAndNonWip() {
	if !d.columnExists("settings", "filter_wip") {
		_, err := d.instance.Exec(`
			ALTER TABLE settings ADD COLUMN filter_wip BOOLEAN DEFAULT 0;
			ALTER TABLE settings ADD COLUMN filter_non_wip BOOLEAN DEFAULT 0;
		`)
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) addSettingsPlannedAndNonPlanned() {
	if !d.columnExists("settings", "planned") {
		_, err := d.instance.Exec(`
			ALTER TABLE settings ADD COLUMN planned BOOLEAN DEFAULT 0;
			ALTER TABLE settings ADD COLUMN non_planned BOOLEAN DEFAULT 0;
		`)
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) FindSettings(settingsId string) (models.Settings, error) {
	var settings models.Settings
	var completedFrom, completedTo string
	var tagsText string

	row := d.instance.QueryRow("SELECT "+SETTINGS_COLUMNS+" FROM settings WHERE id = ?", settingsId)
	err := row.Scan(
		&settings.Id,
		&settings.TasksQuery.FilterCompleted,
		&settings.TasksQuery.FilterIncompleted,
		&settings.TasksQuery.SortColumn,
		&settings.TasksQuery.SortDirection,
		&completedFrom,
		&completedTo,
		&settings.TasksQuery.FilterWip,
		&settings.TasksQuery.FilterNonWip,
		&settings.TasksQuery.Planned,
		&settings.TasksQuery.NonPlanned,
		&tagsText,
		&settings.TasksQuery.SearchText,
		&settings.TasksQuery.EnableLimit,
		&settings.TasksQuery.LimitCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Settings{}, ErrNotFound
		}
		return models.Settings{}, fmt.Errorf("failed to fetch settings: %s: %w", settingsId, err)
	}

	if completedFrom == "" {
		completedFrom = time.Time{}.Format(consts.DEFAULT_DATE_FORMAT)
	}
	settings.TasksQuery.CompletedFrom, err = time.Parse(consts.DEFAULT_DATE_FORMAT, completedFrom)
	if err != nil {
		return models.Settings{}, fmt.Errorf("failed to parse completed_from: %w", err)
	}
	if completedTo == "" {
		completedTo = time.Time{}.Format(consts.DEFAULT_DATE_FORMAT)
	}
	settings.TasksQuery.CompletedTo, err = time.Parse(consts.DEFAULT_DATE_FORMAT, completedTo)
	if err != nil {
		return models.Settings{}, fmt.Errorf("failed to parse completed_to: %w", err)
	}
	if tagsText != "" {
		err = json.Unmarshal([]byte(tagsText), &settings.TasksQuery.Tags)
		if err != nil {
			return models.Settings{}, fmt.Errorf("failed to pars tags: %w", err)
		}
	}
	return settings, nil
}

func (d *DbSQLite) SaveSettings(s models.Settings) error {
	sqlQuery :=
		"INSERT INTO settings (" + SETTINGS_COLUMNS + ") " +
			`VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            filter_completed=excluded.filter_completed,
            filter_incompleted=excluded.filter_incompleted,
            active_sort_column=excluded.active_sort_column,
            active_sort_direction=excluded.active_sort_direction,
            completed_from=excluded.completed_from,
            completed_to=excluded.completed_to,
            filter_wip=excluded.filter_wip,
            filter_non_wip=excluded.filter_non_wip,
            planned=excluded.planned,
            non_planned=excluded.non_planned,
			tags=excluded.tags,
			search_text=excluded.search_text,
			enable_limit=excluded.enable_limit,
			limit_count=excluded.limit_count
    `
	completedFrom := time.Time{}.Format(consts.DEFAULT_DATE_FORMAT)
	if !s.TasksQuery.CompletedFrom.IsZero() {
		completedFrom = s.TasksQuery.CompletedFrom.Format(consts.DEFAULT_DATE_FORMAT)
	}
	completedTo := time.Time{}.Format(consts.DEFAULT_DATE_FORMAT)
	if !s.TasksQuery.CompletedTo.IsZero() {
		completedTo = s.TasksQuery.CompletedTo.Format(consts.DEFAULT_DATE_FORMAT)
	}
	tagsText, err := json.Marshal(s.TasksQuery.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %v: %w", s.TasksQuery.Tags, err)
	}

	common.Debug("SaveSettings: %v", s.TasksQuery)

	args := []any{
		s.Id,
		s.TasksQuery.FilterCompleted,
		s.TasksQuery.FilterIncompleted,
		s.TasksQuery.SortColumn,
		s.TasksQuery.SortDirection,
		completedFrom,
		completedTo,
		s.TasksQuery.FilterWip,
		s.TasksQuery.FilterNonWip,
		s.TasksQuery.Planned,
		s.TasksQuery.NonPlanned,
		tagsText,
		s.TasksQuery.SearchText,
		s.TasksQuery.EnableLimit,
		s.TasksQuery.LimitCount,
	}

	_, err = d.instance.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to save settings: %v: %w", s, err)
	}
	return nil
}
