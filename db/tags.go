package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/models"
	_ "modernc.org/sqlite"
)

const (
	TASKS_TAGS_COLUMNS = "task_id, tag_id"
	TAGS_COLUMNS       = "id, created"
)

func (d *DbSQLite) initTags() {
	common.Debug("initTags")
	d.addTagsTable()
}

func (d *DbSQLite) addTagsTable() {
	id := "add_tags_support"
	if !d.MigrationExists(id) {
		tagsTableSql := `
		CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			created TEXT
		)
		`
		_, err := d.instance.Exec(tagsTableSql)
		if err != nil {
			panic(err)
		}

		tasksTagsSql := `
		CREATE TABLE IF NOT EXISTS TasksTags (
			task_id TEXT,
			tag_id TEXT,
			PRIMARY KEY (task_id, tag_id),
			FOREIGN KEY (task_id) REFERENCES tasks(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)
		`
		_, err = d.instance.Exec(tasksTagsSql)
		if err != nil {
			panic(err)
		}

		d.RecordMigration(id)
	}
}

func (d *DbSQLite) SaveTag(tagId string) error {
	sql := "INSERT INTO tags (" + TAGS_COLUMNS + ") " + " VALUES (?, ?)"
	args := []any{
		tagId,
		time.Now().Format(consts.DEFAULT_TIME_FORMAT),
	}
	logQuery("SaveTag", sql, args)

	_, err := d.instance.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("SaveTag: error; tagId=%v; %w", tagId, err)
	}
	return nil
}

func (d *DbSQLite) AddTagToTask(taskId, tagId string) error {
	sql := "INSERT INTO TasksTags (" + TASKS_TAGS_COLUMNS + ") " + " VALUES (?, ?)"
	args := []any{
		taskId,
		tagId,
	}
	logQuery("AddTagToTask", sql, args)

	_, err := d.instance.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to add tag to task; taskId=%v; tagId=%v; %w", taskId, tagId, err)
	}
	return nil
}

func (d *DbSQLite) deleteAllTagsFromTask(taskId string, tx *sql.Tx) error {
	sql := "DELETE FROM TasksTags WHERE task_id = ?"
	args := []any{
		taskId,
	}
	logQuery("deleteAllTagsFromTask", sql, args)
	_, err := tx.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("deleteAllTagsFromTask: %w", err)
	}
	return nil
}

func (d *DbSQLite) DeleteTagFromTask(taskId, tagId string) error {
	sql := "DELETE FROM TasksTags WHERE task_id = ? AND tag_id = ?"
	args := []any{
		taskId,
		tagId,
	}
	logQuery("DeleteTagFromTask", sql, args)

	result, err := d.instance.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteTagFromTask: failed to delete tag from task; taskId=%v; tagId=%v; %w", taskId, tagId, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteTagFromTask: failed to get affected rows; taskId=%v; tagId=%v; %w", taskId, tagId, err)
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

func (d *DbSQLite) TaskTags(taskId string) ([]models.TaskTag, error) {
	sql := "SELECT " + TASKS_TAGS_COLUMNS + " FROM TasksTags WHERE task_id = ?"
	args := []interface{}{taskId}
	logQuery("TaskTags", sql, args)

	rows, err := d.instance.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TaskTags: failed to query tags for task %s: %w", taskId, err)
	}
	defer rows.Close()

	var tags []models.TaskTag
	for rows.Next() {
		var tagId, taskId string
		if err := rows.Scan(&taskId, &tagId); err != nil {
			return nil, fmt.Errorf("TaskTags: failed to scan tag: %w", err)
		}
		tags = append(tags, models.TaskTag(tagId))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("TaskTags: error iterating tags: %w", err)
	}

	return tags, nil
}

func (d *DbSQLite) TasksTags(taskIds []string) (map[string][]models.TaskTag, error) {
	result := make(map[string][]models.TaskTag)
	if len(taskIds) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(taskIds))
	args := make([]any, len(taskIds))
	for i, id := range taskIds {
		placeholders[i] = "?"
		args[i] = id
	}

	sql := fmt.Sprintf("SELECT %s FROM TasksTags WHERE task_id IN (%s)",
		TASKS_TAGS_COLUMNS,
		strings.Join(placeholders, ","))

	logQuery("TasksTags", sql, args)

	rows, err := d.instance.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TasksTags: failed to query tags for tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var taskId, tagId string
		if err := rows.Scan(&taskId, &tagId); err != nil {
			return nil, fmt.Errorf("TasksTags: failed to scan tag: %w", err)
		}
		result[taskId] = append(result[taskId], models.TaskTag(tagId))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("TasksTags: error iterating tags: %w", err)
	}

	return result, nil
}

func (d *DbSQLite) Tags() ([]models.TaskTag, error) {
	sql := "SELECT id FROM tags ORDER BY created DESC"
	logQuery("Tags", sql, nil)

	rows, err := d.instance.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("Tags: failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []models.TaskTag
	for rows.Next() {
		var tagId string
		if err := rows.Scan(&tagId); err != nil {
			return nil, fmt.Errorf("Tags: failed to scan tag: %w", err)
		}
		tags = append(tags, models.TaskTag(tagId))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Tags: error iterating tags: %w", err)
	}

	return tags, nil
}
