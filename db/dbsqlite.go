package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/models"
	_ "modernc.org/sqlite"
)

type DbSQLite struct {
	instance *sql.DB
}

func NewDbSQLite() *DbSQLite {
	return &DbSQLite{}
}

const (
	SETTINGS_COLUMNS = "id, filter_completed, filter_incompleted, active_sort_column, active_sort_direction, completed_from, completed_to"
)

func (d *DbSQLite) Init() {
	dir, err := common.ResolveAppDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	file := filepath.Join(dir, "db.sqlite")

	db, err := sql.Open("sqlite", file)
	if err != nil {
		log.Fatal(err)
	}
	d.instance = db

	d.initTasks()

	_, err = d.instance.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id TEXT PRIMARY KEY,
			filter_completed BOOLEAN,
			filter_incompleted BOOLEAN,
			active_sort_column INTEGER,
			active_sort_direction INTEGER,
			completed_from TEXT,
			completed_to TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	d.addSettingsCompletedFrom()
	d.addSettingsCompletedTo()
	d.addSettingsFilterIncomplete()
}

func (d *DbSQLite) columnExists(tableName, columnName string) bool {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	rows, err := d.instance.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype, notnull string
		var dfltValue interface{}
		var primaryKey int

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &primaryKey); err != nil {
			panic(err)
		}

		if name == columnName {
			return true
		}
	}
	return false
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

func (d *DbSQLite) Close() {
	common.Debug("closing db...")
	d.instance.Close()
}

func (d *DbSQLite) FindSettings(settingsId string) (models.Settings, error) {
	var settings models.Settings
	var completedFrom, completedTo string

	row := d.instance.QueryRow("SELECT "+SETTINGS_COLUMNS+" FROM settings WHERE id = ?", settingsId)
	err := row.Scan(
		&settings.Id,
		&settings.TasksQuery.FilterCompleted,
		&settings.TasksQuery.FilterIncompleted,
		&settings.TasksQuery.SortColumn,
		&settings.TasksQuery.SortDirection,
		&completedFrom,
		&completedTo,
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

	return settings, nil
}

func (d *DbSQLite) SaveSettings(s models.Settings) error {
	sqlQuery :=
		"INSERT INTO settings (" + SETTINGS_COLUMNS + ") " +
			`VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            filter_completed=excluded.filter_completed,
            filter_incompleted=excluded.filter_incompleted,
            active_sort_column=excluded.active_sort_column,
            active_sort_direction=excluded.active_sort_direction,
            completed_from=excluded.completed_from,
            completed_to=excluded.completed_to
    `
	completedFrom := time.Time{}.Format(consts.DEFAULT_DATE_FORMAT)
	if !s.TasksQuery.CompletedFrom.IsZero() {
		completedFrom = s.TasksQuery.CompletedFrom.Format(consts.DEFAULT_DATE_FORMAT)
	}

	completedTo := time.Time{}.Format(consts.DEFAULT_DATE_FORMAT)
	if !s.TasksQuery.CompletedTo.IsZero() {
		completedTo = s.TasksQuery.CompletedTo.Format(consts.DEFAULT_DATE_FORMAT)
	}

	args := []interface{}{
		s.Id,
		s.TasksQuery.FilterCompleted,
		s.TasksQuery.FilterIncompleted,
		s.TasksQuery.SortColumn,
		s.TasksQuery.SortDirection,
		completedFrom,
		completedTo,
	}

	_, err := d.instance.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to save settings: %v: %w", s, err)
	}
	return nil
}
