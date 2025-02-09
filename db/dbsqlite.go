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

	_, err = d.instance.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT,
			content TEXT,
			created TEXT,
			updated TEXT,
			completed TEXT,
			priority INTEGER,
			wip INTEGER DEFAULT 0,
			planned INTEGER DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS settings (
			id TEXT PRIMARY KEY,
			filter_completed BOOLEAN,
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
	d.addTasksWipColumn()
	d.addTasksPlannedColumn()
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

func (d *DbSQLite) addTasksWipColumn() {
	if !d.columnExists("tasks", "wip") {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN wip INTEGER DEFAULT 0")
		if err != nil {
			panic(err)
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

func (d *DbSQLite) addTasksPlannedColumn() {
	if !d.columnExists("tasks", "planned") {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN planned INTEGER DEFAULT 0")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) Close() {
	common.Debug("closing db...")
	d.instance.Close()
}

func (d *DbSQLite) scanNextTask(rows *sql.Rows) (models.Task, error) {
	var task models.Task
	var created, updated, completed string
	var wip, planned int

	err := rows.Scan(&task.Id, &task.Title, &task.Content, &created, &updated, &completed, &task.Priority, &wip, &planned)
	if err != nil {
		return models.EMPTY_TASK, err
	}

	task.Wip = wip == 1
	task.Planned = planned == 1

	task.Created, err = time.Parse(consts.DEFAULT_TIME_FORMAT, created)
	if err != nil {
		return models.EMPTY_TASK, fmt.Errorf("failed to parse created time: %w", err)
	}

	task.Updated, err = time.Parse(consts.DEFAULT_TIME_FORMAT, updated)
	if err != nil {
		return models.EMPTY_TASK, fmt.Errorf("failed to parse updated time: %w", err)
	}

	task.Completed, err = time.Parse(consts.DEFAULT_TIME_FORMAT, completed)
	if err != nil {
		return models.EMPTY_TASK, fmt.Errorf("failed to parse completed time: %w", err)
	}

	return task, nil
}

func (d *DbSQLite) Tasks() (result []models.Task, err error) {
	rows, err := d.instance.Query("SELECT id, title, content, created, updated, completed, priority, wip, planned FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task, err := d.scanNextTask(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, task)
	}
	return result, nil
}

func (d *DbSQLite) FindTask(taskId string) (models.Task, error) {
	rows, err := d.instance.Query("SELECT id, title, content, created, updated, completed, priority, wip, planned FROM tasks WHERE id = ?", taskId)
	if err != nil {
		return models.EMPTY_TASK, fmt.Errorf("failed to query task: %s: %w", taskId, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return models.EMPTY_TASK, ErrNotFound
	}

	return d.scanNextTask(rows)
}

func (d *DbSQLite) DeleteTask(taskId string) error {
	result, err := d.instance.Exec("DELETE FROM tasks WHERE id = ?", taskId)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	rc, err := result.RowsAffected()
	if rc == 0 || err != nil {
		return ErrNotFound
	}

	return nil
}

func (d *DbSQLite) DeleteAllTasks() error {
	_, err := d.instance.Exec("DELETE FROM tasks")
	if err != nil {
		return fmt.Errorf("failed to delete all tasks: %v", err)
	}
	return nil
}

func (d *DbSQLite) SaveTask(task models.Task) error {
	_, err := d.instance.Exec(`
		INSERT INTO tasks (id, title, content, created, updated, completed, priority, wip, planned)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title,
			content=excluded.content,
			created=excluded.created,
			updated=excluded.updated,
			completed=excluded.completed,
			priority=excluded.priority,
			wip=excluded.wip,
			planned=excluded.planned
	`,
		task.Id,
		task.Title,
		task.Content,
		task.Created.Format(consts.DEFAULT_TIME_FORMAT),
		task.Updated.Format(consts.DEFAULT_TIME_FORMAT),
		task.Completed.Format(consts.DEFAULT_TIME_FORMAT),
		task.Priority,
		task.Wip,
		task.Planned,
	)
	if err != nil {
		return fmt.Errorf("failed to save task: %v: %w", task, err)
	}
	return nil
}

func (d *DbSQLite) FindSettings(settingsId string) (models.Settings, error) {
	var settings models.Settings
	var completedFrom, completedTo string

	row := d.instance.QueryRow(`
        SELECT id, filter_completed, active_sort_column, active_sort_direction, 
               completed_from, completed_to 
        FROM settings 
        WHERE id = ?`, settingsId)

	err := row.Scan(
		&settings.Id,
		&settings.TasksQuery.FilterCompleted,
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
	sqlQuery := `
        INSERT INTO settings (
            id, 
            filter_completed, 
            active_sort_column, 
            active_sort_direction,
            completed_from,
            completed_to
        )
        VALUES (?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            filter_completed=excluded.filter_completed,
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

func (d *DbSQLite) FindTasks(query models.TasksQuery) ([]models.Task, error) {
	var args []interface{}
	sqlQuery := "SELECT id, title, content, created, updated, completed, priority, wip, planned FROM tasks WHERE 1=1"

	common.Debug("FindTasks: query: %v", query)

	if query.FilterCompleted {
		sqlQuery += " AND completed = ?"
		notCompleted := models.NOT_COMPLETED.Format(consts.DEFAULT_TIME_FORMAT)
		args = append(args, notCompleted)
	} else {
		if !query.CompletedFrom.IsZero() {
			sqlQuery += " AND (completed >= ? OR completed = ?)"
			args = append(args, query.CompletedFrom.Format(consts.DEFAULT_TIME_FORMAT), models.NOT_COMPLETED.Format(consts.DEFAULT_TIME_FORMAT))
		}

		if !query.CompletedTo.IsZero() {
			sqlQuery += " AND (completed <= ? OR completed = ?)"
			args = append(args, query.CompletedTo.Format(consts.DEFAULT_TIME_FORMAT), models.NOT_COMPLETED.Format(consts.DEFAULT_TIME_FORMAT))
		}
	}

	if query.SortColumn != models.ColumnUndefined {
		sqlQuery += " ORDER BY "
		switch query.SortColumn {
		case models.Completed:
			sqlQuery += "completed"
		case models.Created:
			sqlQuery += "created"
		case models.Priority:
			sqlQuery += "priority"
		default:
			sqlQuery += "created" // default sort
		}

		if query.SortDirection == models.Desc {
			sqlQuery += " DESC"
		} else {
			sqlQuery += " ASC"
		}
	}

	common.Debug("FindTasks: sqlQuery: %v", sqlQuery)
	common.Debug("FindTasks: args: %v", args)

	rows, err := d.instance.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var result []models.Task
	for rows.Next() {
		task, err := d.scanNextTask(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return result, nil
}
