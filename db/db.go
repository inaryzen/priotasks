package db

import (
	"errors"

	"github.com/inaryzen/priotasks/models"
)

var ErrNotFound = errors.New("not found")

const (
	PREFIX_TASK     = "card:"
	PREFIX_SETTINGS = "settings:"
)

var instance Db

type Db interface {
	Init(string)
	Close()
	Tasks() ([]models.Task, error)
	FindTask(taskId string) (models.Task, error)
	FindTasks(query models.TasksQuery) ([]models.Task, error)
	DeleteTask(taskId string) error
	DeleteAllTasks() error
	SaveTask(task models.Task) error
	FindSettings(settingsId string) (models.Settings, error)
	SaveSettings(s models.Settings) error
	MigrationExists(id string) bool
	RecordMigration(id string)
	SaveTag(tagId string) error
	AddTagToTask(taskId, tagId string) error
	TaskTags(taskId string) ([]models.TaskTag, error)
	Tags() ([]models.TaskTag, error)
}

func SetDB(db Db) {
	instance = db
}

func DB() Db {
	return instance
}
