package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/models"
	_ "modernc.org/sqlite"
)

const (
	TASK_COLUMNS = "id, title, content, created, updated, completed, priority, wip, planned, impact, cost, value, fun"
)

func (d *DbSQLite) initTasks() {
	common.Debug("initTasks")
	var err error
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
			planned INTEGER DEFAULT 0,
			impact INTEGER DEFAULT 2,
			cost INTEGER DEFAULT 2
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	d.addTasksWipColumn()
	d.addTasksPlannedColumn()
	d.addTasksImpactColumn()
	d.addTasksCostColumn()
	d.addValueColumn()
	d.addTasksFunColumn()
}

func (d *DbSQLite) addValueColumn() {
	id := "task_table_add_value_column"
	if !d.MigrationExists(id) {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN value REAL DEFAULT 0")
		if err != nil {
			panic(err)
		} else {
			d.RecordMigration(id)
		}
	}
}

func (d *DbSQLite) addTasksWipColumn() {
	if !d.columnExists("tasks", "wip") {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN wip INTEGER DEFAULT 0")
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

func (d *DbSQLite) addTasksCostColumn() {
	if !d.columnExists("tasks", "cost") {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN cost INTEGER DEFAULT 2")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) addTasksImpactColumn() {
	if !d.columnExists("tasks", "impact") {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN impact INTEGER DEFAULT 2")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) addTasksFunColumn() {
	if !d.columnExists("tasks", "fun") {
		_, err := d.instance.Exec("ALTER TABLE tasks ADD COLUMN fun INTEGER DEFAULT 1")
		if err != nil {
			panic(err)
		}
	}
}

func (d *DbSQLite) scanNextTask(rows *sql.Rows) (models.Task, error) {
	var task models.Task
	var created, updated, completed string
	var wip, planned int

	err := rows.Scan(
		&task.Id,
		&task.Title,
		&task.Content,
		&created,
		&updated,
		&completed,
		&task.Priority,
		&wip,
		&planned,
		&task.Impact,
		&task.Cost,
		&task.Value,
		&task.Fun,
	)
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
	rows, err := d.instance.Query("SELECT " + TASK_COLUMNS + " FROM tasks")
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
	rows, err := d.instance.Query("SELECT "+TASK_COLUMNS+" FROM tasks WHERE id = ?", taskId)
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

	tx, err := d.instance.Begin()
	if err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
	}

	err = d.deleteAllTagsFromTask(taskId, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteTask: %w", err)
	}

	result, err := tx.Exec("DELETE FROM tasks WHERE id = ?", taskId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete task: %v", err)
	}

	rc, err := result.RowsAffected()
	if rc == 0 || err != nil {
		tx.Rollback()
		return ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
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
	sql := "INSERT INTO tasks (" + TASK_COLUMNS + ") " +
		`VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title=excluded.title,
			content=excluded.content,
			created=excluded.created,
			updated=excluded.updated,
			completed=excluded.completed,
			priority=excluded.priority,
			wip=excluded.wip,
			planned=excluded.planned,
			impact=excluded.impact,
			cost=excluded.cost,
			value=excluded.value,
			fun=excluded.fun
	`
	args := []any{
		task.Id,
		task.Title,
		task.Content,
		task.Created.Format(consts.DEFAULT_TIME_FORMAT),
		task.Updated.Format(consts.DEFAULT_TIME_FORMAT),
		task.Completed.Format(consts.DEFAULT_TIME_FORMAT),
		task.Priority,
		task.Wip,
		task.Planned,
		task.Impact,
		task.Cost,
		task.Value,
		task.Fun,
	}
	logQuery("SaveTask", sql, args)
	_, err := d.instance.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to save task: %v: %w", task, err)
	}
	return nil
}

func logQuery(prefix, sql string, args []interface{}) {
	common.Debug("%v: sqlQuery: %v", prefix, sql)
	common.Debug("%v: args: %v", prefix, args)
}

func (d *DbSQLite) FindTasks(query models.TasksQuery) ([]models.Task, error) {
	var args []any
	sqlQuery := "SELECT " + TASK_COLUMNS + " FROM tasks"
	sqlQuery += " WHERE 1=1"

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
	if query.FilterIncompleted {
		sqlQuery += " AND completed != ?"
		notCompleted := models.NOT_COMPLETED.Format(consts.DEFAULT_TIME_FORMAT)
		args = append(args, notCompleted)
	}
	if query.FilterWip {
		sqlQuery += " AND wip = 1"
	}
	if query.FilterNonWip {
		sqlQuery += " AND wip = 0"
	}
	if query.Planned {
		sqlQuery += " AND planned = 1"
	}
	if query.NonPlanned {
		sqlQuery += " AND planned = 0"
	}
	if len(query.Tags) > 0 {
		sqlQuery += " AND id in ( select task_id from TasksTags where tag_id in ("
		for i, t := range query.Tags {
			if i != 0 {
				sqlQuery += ", "
			}
			sqlQuery = sqlQuery + "?"
			args = append(args, t)
		}
		sqlQuery += "))"
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
		case models.ColumnImpact:
			sqlQuery += "impact"
		case models.ColumnWip:
			sqlQuery += "wip"
		case models.ColumnPlanned:
			sqlQuery += "planned"
		case models.ColumnCost:
			sqlQuery += "cost"
		case models.ColumnValue:
			sqlQuery += "value"
		case models.ColumnFun:
			sqlQuery += "fun"
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
