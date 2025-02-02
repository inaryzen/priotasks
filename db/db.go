package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/models"
)

var dbBadger *badger.DB
var ErrNotFound = errors.New("not found")

const (
	PREFIX_TASK     = "card:"
	PREFIX_SETTINGS = "settings:"
)

func DbInit() {
	dir, err := common.ResolveAppDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	dir = filepath.Join(dir, "db")

	d, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		log.Fatal(err)
	}
	dbBadger = d
}

func DbClose() {
	dbBadger.Close()
}

func Tasks() (result []models.Task, err error) {

	result = []models.Task{}
	err = nil

	err = dbBadger.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PREFIX_TASK)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()

			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				log.Fatalf("Error retrieving value: %v: %v", key, err)
			}

			var task models.Task
			err = json.Unmarshal(valueCopy, &task)
			if err != nil {
				log.Fatalf("Error unmarshalling value: %v: %v", key, err)
			}

			result = append(result, task)
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("failed to fetch records: %w", err)
	}
	return
}

// TODO: handle the case when Card is not found properly, see ErrKeyNotFound
func FindTask(taskId string) (models.Task, error) {
	var result models.Task
	var err error

	err = dbBadger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(PREFIX_TASK + taskId))

		errors.Is(err, badger.ErrKeyNotFound)
		if err != nil {
			return errors.Join(ErrNotFound, err)
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(value, &result)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		err = fmt.Errorf("failed to fetch records: %s: %w", taskId, err)
	}
	return result, err
}

func DeleteAllTasks() error {
	err := dbBadger.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PREFIX_TASK)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func SaveTask(task models.Task) error {
	var err error

	jsonData, err := json.Marshal(task)
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	err = dbBadger.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(PREFIX_TASK+task.Id), jsonData)
		return err
	})
	return logSavedEntity(task, err)
}

func FindSettings(settingsId string) (models.Settings, error) {
	var result models.Settings
	var err error

	err = dbBadger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(PREFIX_SETTINGS + settingsId))

		errors.Is(err, badger.ErrKeyNotFound)
		if err != nil {
			return errors.Join(ErrNotFound, err)
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(value, &result)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		err = fmt.Errorf("failed to fetch records: %s: %w", settingsId, err)
	}
	return result, err
}

func SaveSettings(s models.Settings) error {
	var err error

	jsonData, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	err = dbBadger.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(PREFIX_SETTINGS+s.Id), jsonData)
		return err
	})

	return logSavedEntity(s, err)
}

func logSavedEntity(entity any, err error) error {
	if err != nil {
		return fmt.Errorf("failed to save entity: %v: %w", entity, err)
	} else {
		if common.IsDebug() {
			log.Printf("Saved: %v \n", entity)
		}
		return nil
	}
}
