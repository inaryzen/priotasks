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

func UpdateTask(changed models.Task) error {
	orig, err := db.DB().FindTask(changed.Id)
	if err != nil {
		log.Printf("failed to find the record: %s: %s", changed.Id, err)
		return err
	}
	orig = orig.Update(changed)
	if err = SaveTask(orig); err != nil {
		return err
	}
	return nil
}

func SaveNewTask(t models.Task) error {
	t = t.AsNewTask()
	if err := SaveTask(t); err != nil {
		return err
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
