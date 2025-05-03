package db

import (
	"github.com/inaryzen/priotasks/models"
)

type NoOpDB struct{}

func (m *NoOpDB) Init(p string)                                            {}
func (m *NoOpDB) Close()                                                   {}
func (m *NoOpDB) Tasks() ([]models.Task, error)                            { return nil, nil }
func (m *NoOpDB) FindTasks(query models.TasksQuery) ([]models.Task, error) { return nil, nil }
func (m *NoOpDB) DeleteTask(taskId string) error                           { return nil }
func (m *NoOpDB) DeleteAllTasks() error                                    { return nil }
func (m *NoOpDB) FindSettings(settingsId string) (models.Settings, error) {
	return models.Settings{}, nil
}
func (m *NoOpDB) SaveSettings(s models.Settings) error                            { return nil }
func (m *NoOpDB) SaveTag(tagId string) error                                      { return nil }
func (m *NoOpDB) TaskTags(taskId string) ([]models.TaskTag, error)                { return nil, nil }
func (m *NoOpDB) Tags() ([]models.TaskTag, error)                                 { return nil, nil }
func (m *NoOpDB) AddTagToTask(taskId, tagId string) error                         { return nil }
func (m *NoOpDB) DeleteTagFromTask(taskId, tagId string) error                    { return nil }
func (m *NoOpDB) FindTask(taskId string) (models.Task, error)                     { return models.Task{}, nil }
func (m *NoOpDB) SaveTask(task models.Task) error                                 { return nil }
func (m *NoOpDB) MigrationExists(id string) bool                                  { return false }
func (m *NoOpDB) RecordMigration(id string)                                       {}
func (m *NoOpDB) TasksTags(taskIds []string) (map[string][]models.TaskTag, error) { return nil, nil }
func (m *NoOpDB) DeleteTag(tagId string) error                                    { return nil }
func (m *NoOpDB) DeleteTagFromAllTasks(tagId string) error                        { return nil }
