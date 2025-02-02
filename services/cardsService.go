package services

import (
	"cmp"
	"fmt"
	"log"
	"slices"

	"github.com/inaryzen/prio_cards/db"
	"github.com/inaryzen/prio_cards/models"
)

func FindCards(FilterCompleted bool, sort models.SortColumn, dir models.SortDirection) ([]models.Card, error) {
	var cards []models.Card
	var err error

	cards, err = db.Cards()
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

	slices.SortFunc(cards, func(a, b models.Card) int {
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

func UpdateCard(c models.Card) error {
	card, err := db.FindCard(c.Id)
	if err != nil {
		log.Printf("failed to find the record: %s: %s", c.Id, err)
		return err
	}
	card = card.Update(c)
	err = db.SaveCard(card)
	if err != nil {
		log.Printf("failed to update the record: %s: %s", card.Id, err)
		return err
	}
	return nil
}

func SaveCard(c models.Card) error {
	err := db.SaveCard(c)
	if err != nil {
		log.Printf("failed to update the record: %s: %s", c.Id, err)
		return err
	}
	return nil
}

func FlipCard(card models.Card) error {
	if card.Completed == models.NOT_COMPLETED {
		card = card.Complete()
	} else {
		card = card.Uncomplete()
	}
	err := db.SaveCard(card)
	if err != nil {
		return fmt.Errorf("failed to flip the card: %v: %w", card.Id, err)
	}
	log.Printf("Updated Completed status of card: %v", card)
	return nil
}
