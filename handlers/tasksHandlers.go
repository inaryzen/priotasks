package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/inaryzen/priotasks/components"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/db"
	"github.com/inaryzen/priotasks/models"
	"github.com/inaryzen/priotasks/services"
)

func resolveTaskOrNotFound(w http.ResponseWriter, r *http.Request) (models.Task, error) {
	idString := r.PathValue("id")
	card, err := db.DB().FindTask(idString)
	if errors.Is(err, db.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	return card, err
}

func DeleteTasksId(w http.ResponseWriter, r *http.Request) {
	card, err := resolveTaskOrNotFound(w, r)
	if err != nil {
		return
	}
	err = services.DeleteTask(card.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	drawTaskTable(w, r)
}

func PostTaskToggleCompleted(w http.ResponseWriter, r *http.Request) {
	card, err := resolveTaskOrNotFound(w, r)
	if err != nil {
		return
	}
	err = services.FlipTask(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func GetViewEmptyTask(w http.ResponseWriter, r *http.Request) {
	cardsView := components.ModalTaskView(models.EMPTY_TASK)
	cardsView.Render(r.Context(), w)
}

func GetViewTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	card, err := resolveTaskOrNotFound(w, r)
	if err != nil {
		return
	}

	cardsView := components.ModalTaskView(card)
	cardsView.Render(r.Context(), w)
}

// TODO: Optimize. Instead of reloading the whole table, update only a single row
func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	c := resolveTaskFromForm(r)
	err := services.UpdateTask(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	drawTaskTable(w, r)
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	c := resolveTaskFromForm(r)
	card := models.Create(c)
	err := services.SaveTask(card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	drawTaskTable(w, r)
}

func resolveTaskFromForm(r *http.Request) models.Task {
	formPriority := r.FormValue("modal-task-priority")
	prio, err := models.StrToTaskPriority(formPriority)
	if err != nil {
		log.Printf("failed to parse Priority: %v: %v", formPriority, err)
	}

	formImpact := r.FormValue("modal-task-impact")
	impact, err := models.StrToImpact(formImpact)
	if err != nil {
		log.Printf("failed to parse Impact: %v: %v", formImpact, err)
	}

	// Parse checkbox values - they will be "on" if checked, or empty if unchecked
	wipValue := r.FormValue("task-wip") == "on"
	plannedValue := r.FormValue("task-planned") == "on"

	return models.Task{
		Id:       r.FormValue("card-id"),
		Content:  r.FormValue("card-text"),
		Title:    r.FormValue("card-title"),
		Priority: prio,
		Wip:      wipValue,
		Planned:  plannedValue,
		Impact:   impact,
	}
}

func drawTaskTable(w http.ResponseWriter, r *http.Request) {
	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	cards, err := findTasksOrWriteError(w)
	if err != nil {
		return
	}
	cardsView := components.TaskTable(cards, settings)
	cardsView.Render(r.Context(), w)
}

func findTasksOrWriteError(w http.ResponseWriter) (cards []models.Task, err error) {
	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	cards, err = services.FindTasks(settings.TasksQuery)

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

func GetTasks(w http.ResponseWriter, r *http.Request) {
	cards, err := findTasksOrWriteError(w)
	if err != nil {
		return
	}

	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	cardsView := components.TasksView(cards, settings)
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

	drawTaskTable(w, r)
}
