package services

import (
	"fmt"

	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

func Tags() ([]models.TaskTag, error) {
	tags, err := db.DB().Tags()
	if err != nil {
		return nil, fmt.Errorf("Tags: %w", err)
	} else {
		return tags, nil
	}
}

func SaveTag(tag models.TaskTag) error {
	err := db.DB().SaveTag(string(tag))
	if err != nil {
		return fmt.Errorf("SaveTag: error tag=%v: %w", tag, err)
	} else {
		return nil
	}
}

func AddTagToTask(taskId string, tag models.TaskTag) error {
	err := db.DB().AddTagToTask(taskId, string(tag))
	if err != nil {
		return fmt.Errorf("AddTagToTask: error tag=%v; taskId=%v: %w", tag, taskId, err)
	} else {
		return nil
	}
}

func DeleteTag(tag models.TaskTag) error {
	return nil
}

func TaskTags(taskId string) ([]models.TaskTag, error) {
	tags, err := db.DB().TaskTags(taskId)
	if err != nil {
		return nil, fmt.Errorf("TaskTags: failed to get tags for task %s: %w", taskId, err)
	}
	return tags, nil
}

func RemoveTagFromTask(taskId string, tag models.TaskTag) error {
	err := db.DB().DeleteTagFromTask(taskId, string(tag))
	if err != nil {
		return fmt.Errorf("RemoveTagFromTask: taskId=%v, tag=%v: %w", taskId, tag, err)
	}
	return nil
}

func TasksTags(taskIds []string) (map[string][]models.TaskTag, error) {
	tags, err := db.DB().TasksTags(taskIds)
	if err != nil {
		return nil, fmt.Errorf("TasksTags: failed to get tags for tasks: %w", err)
	}
	return tags, nil
}
