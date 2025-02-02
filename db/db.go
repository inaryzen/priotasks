package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"

	"github.com/inaryzen/prio_cards/common"
	"github.com/inaryzen/prio_cards/models"
)

var dbBadger *badger.DB
var ErrNotFound = errors.New("not found")

const (
	PREFIX_CARD     = "card:"
	PREFIX_SETTINGS = "settings:"
)

func DbInit() {
	log.Printf("Opening Badger...")
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
	log.Printf("Closing Badger...")
	dbBadger.Close()
}

func Cards() (result []models.Card, err error) {

	result = []models.Card{}
	err = nil

	// return cards
	err = dbBadger.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PREFIX_CARD)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()

			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				log.Fatalf("Error retrieving value: %v: %v", key, err)
			}

			var card models.Card
			err = json.Unmarshal(valueCopy, &card)
			if err != nil {
				log.Fatalf("Error unmarshalling value: %v: %v", key, err)
			}

			result = append(result, card)
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("failed to fetch records: %w", err)
	}
	return
}

// TODO: handle the case when Card is not found properly, see ErrKeyNotFound
func FindCard(cardId string) (models.Card, error) {
	var result models.Card
	var err error

	// return cards
	err = dbBadger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(PREFIX_CARD + cardId))

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
		err = fmt.Errorf("failed to fetch records: %s: %w", cardId, err)
	}
	return result, err
}

func DeleteAllCards() error {
	err := dbBadger.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(PREFIX_CARD)

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

func SaveCard(card models.Card) error {
	var err error

	jsonData, err := json.Marshal(card)
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	err = dbBadger.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(PREFIX_CARD+card.Id), jsonData)
		return err
	})
	return logSavedEntity(card, err)
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
