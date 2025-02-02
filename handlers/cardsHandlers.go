package handlers

import (
	"errors"
	"log"

	"net/http"

	"github.com/inaryzen/prio_cards/components"
	"github.com/inaryzen/prio_cards/consts"
	"github.com/inaryzen/prio_cards/db"
	"github.com/inaryzen/prio_cards/models"
	"github.com/inaryzen/prio_cards/services"
)

func resolveCardOrNotFound(w http.ResponseWriter, r *http.Request) (models.Card, error) {
	idString := r.PathValue("id")
	card, err := db.FindCard(idString)
	if errors.Is(err, db.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	return card, err
}

func PostCardToggleCompleted(w http.ResponseWriter, r *http.Request) {
	card, err := resolveCardOrNotFound(w, r)
	if err != nil {
		return
	}
	err = services.FlipCard(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func GetViewEmptyCard(w http.ResponseWriter, r *http.Request) {
	cardsView := components.ModalCardView(models.EMPTY_CARD)
	cardsView.Render(r.Context(), w)
}

func GetViewCardByIdHandler(w http.ResponseWriter, r *http.Request) {
	card, err := resolveCardOrNotFound(w, r)
	if err != nil {
		return
	}

	cardsView := components.ModalCardView(card)
	cardsView.Render(r.Context(), w)
}

// TODO: Optimize. Instead of reloading the whole table, update only a single row
func PutCardHandler(w http.ResponseWriter, r *http.Request) {
	c := resolveCardFromForm(r)
	err := services.UpdateCard(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	drawCardTable(w, r)
}

func PostCardHandler(w http.ResponseWriter, r *http.Request) {
	c := resolveCardFromForm(r)
	card := models.Create(c)
	err := services.SaveCard(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	drawCardTable(w, r)
}

func resolveCardFromForm(r *http.Request) models.Card {
	formPriority := r.FormValue("modal-task-priority")
	prio, err := models.StrToTaskPriority(formPriority)
	if err != nil {
		log.Printf("failed to parse Priority: %v: %v", formPriority, err)
	}

	return models.Card{
		Id:       r.FormValue("card-id"),
		Content:  r.FormValue("card-text"),
		Title:    r.FormValue("card-title"),
		Priority: prio,
	}
}

func drawCardTable(w http.ResponseWriter, r *http.Request) {
	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	cards, err := findCardsOrWriteError(w)
	if err != nil {
		return
	}
	cardsView := components.CardTable(cards, settings)
	cardsView.Render(r.Context(), w)
}

func findCardsOrWriteError(w http.ResponseWriter) (cards []models.Card, err error) {
	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	cards, err = services.FindCards(settings.FilterCompleted, settings.ActiveSortColumn, settings.ActiveSortDirection)

	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func findSettingsOrWriteError(w http.ResponseWriter) (models.Settings, error) {
	settings, err := services.FindUserSettings()
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return settings, err
}

func GetCards(w http.ResponseWriter, r *http.Request) {
	cards, err := findCardsOrWriteError(w)
	if err != nil {
		return
	}

	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	cardsView := components.CardsView(cards, settings)
	cardsView.Render(r.Context(), w)
}

func PostToggleSortTable(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}

	var param string
	param = r.Form.Get(consts.SORT_COLUMN_NAME)
	sortColumn := models.ColumnFromString(param)
	param = r.Form.Get(consts.SORT_DIRECTION_NAME)
	sortDirection := models.DirectionFromString(param)

	err = services.ToggleSorting(settings, sortColumn, sortDirection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	drawCardTable(w, r)
}
