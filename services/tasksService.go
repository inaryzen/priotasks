package services

import (
	"cmp"
	"fmt"
	"log"
	"slices"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
)

func FindTasks(FilterCompleted bool, sort models.SortColumn, dir models.SortDirection) ([]models.Task, error) {
	var cards []models.Task
	var err error

	cards, err = db.DB().Tasks()
	if err != nil {
		return cards, fmt.Errorf("failed to retrieve cards: %w", err)
	}

	if FilterCompleted {
		filtered := cards[:0]
		for _, v := range cards {
			if !v.IsCompleted() {
				filtered = append(filtered, v)
			}
		}
		cards = filtered
	}

	slices.SortFunc(cards, func(a, b models.Task) int {
		if sort == models.Completed {
			if dir == models.Desc {
				return cmp.Compare(b.Completed.Unix(), a.Completed.Unix())
			} else {
				return cmp.Compare(a.Completed.Unix(), b.Completed.Unix())
			}
		} else if sort == models.Created {
			if dir == models.Desc {
				return cmp.Compare(b.Created.Unix(), a.Created.Unix())
			} else {
				return cmp.Compare(a.Created.Unix(), b.Created.Unix())
			}
		} else if sort == models.Priority {
			if dir == models.Desc {
				return cmp.Compare(b.Priority, a.Priority)
			} else {
				return cmp.Compare(a.Priority, b.Priority)
			}
		} else {
			return cmp.Compare(b.Created.Unix(), a.Created.Unix())
		}
	})

	return cards, err
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

func UpdateTask(c models.Task) error {
	card, err := db.DB().FindTask(c.Id)
	if err != nil {
		log.Printf("failed to find the record: %s: %s", c.Id, err)
		return err
	}
	card = card.Update(c)
	err = db.DB().SaveTask(card)
	if err != nil {
		log.Printf("failed to update the record: %s: %s", card.Id, err)
		return err
	}
	return nil
}

func SaveTask(c models.Task) error {
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
