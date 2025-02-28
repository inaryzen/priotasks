package services

import (
	"fmt"
	"log"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

func FindTasks(query models.TasksQuery) ([]models.Task, error) {
	tasks, err := db.DB().FindTasks(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}
	return tasks, nil
}

func DeleteTask(taskId string) error {
	common.Debug("deleting task: %v", taskId)

	err := db.DB().DeleteTask(taskId)
	if err != nil {
		log.Printf("failed to delete the task: %s: %s", taskId, err)
		return err
	}
	return nil
}

func DeleteAllTasks() error {
	common.Debug("deleting all tasks")

	err := db.DB().DeleteAllTasks()
	if err != nil {
		log.Printf("failed to delete all tasks: %s", err)
		return err
	}
	return nil
}

func UpdateTask(changed models.Task, changedTags []models.TaskTag) error {
	orig, err := db.DB().FindTask(changed.Id)
	if err != nil {
		log.Printf("failed to find the record: %s: %s", changed.Id, err)
		return err
	}
	orig = orig.Update(changed)
	if err = SaveTask(orig); err != nil {
		return err
	}

	origTags, err := TaskTags(orig.Id)
	if err != nil {
		return fmt.Errorf("UpdateTask: %w", err)
	}

	findMissing := func(source, target []models.TaskTag) []models.TaskTag {
		targetMap := make(map[models.TaskTag]bool)
		for _, item := range target {
			targetMap[item] = true
		}

		var missing []models.TaskTag
		for _, item := range source {
			if !targetMap[item] {
				missing = append(missing, item)
			}
		}
		return missing
	}

	newTags := findMissing(changedTags, origTags)
	for _, t := range newTags {
		err := RemoveTagFromTask(orig.Id, t)
		if err != nil {
			return fmt.Errorf("UpdateTask: %w", err)
		}
	}

	removeTags := findMissing(origTags, changedTags)
	for _, t := range removeTags {
		err := AddTagToTask(orig.Id, t)
		if err != nil {
			return fmt.Errorf("UpdateTask: %w", err)
		}
	}

	return nil
}

func SaveNewTask(t models.Task, tags []models.TaskTag) error {
	t = t.AsNewTask()
	if err := SaveTask(t); err != nil {
		return err
	}

	for _, tag := range tags {
		err := AddTagToTask(t.Id, tag)
		if err != nil {
			return fmt.Errorf("SaveNewTask: %w", err)
		}
	}
	return nil
}

func SaveTask(c models.Task) error {
	c = c.CalculateValue()

	err := db.DB().SaveTask(c)
	if err != nil {
		log.Printf("failed to update the record: %s: %s", c.Id, err)
		return err
	}
	return nil
}

func FlipTask(card models.Task) error {
	if card.Completed == models.NOT_COMPLETED {
		card = card.Complete()
	} else {
		card = card.Uncomplete()
	}
	err := db.DB().SaveTask(card)
	if err != nil {
		return fmt.Errorf("failed to flip the card: %v: %w", card.Id, err)
	}
	if common.IsDebug() {
		log.Printf("Updated Completed status of card: %v", card)
	}
	return nil
}

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
	return nil
}
