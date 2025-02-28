package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

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
	allTags, err := services.Tags()
	if err != nil {
		internalServerError(w, err)
		return
	}

	cardsView := components.ModalTaskView(models.EMPTY_TASK, nil, allTags)
	cardsView.Render(r.Context(), w)
}

func GetViewTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	task, err := resolveTaskOrNotFound(w, r)
	if err != nil {
		internalServerError(w, err)
		return
	}
	allTags, err := services.Tags()
	if err != nil {
		internalServerError(w, err)
		return
	}
	taskTags, err := services.TaskTags(task.Id)
	if err != nil {
		internalServerError(w, err)
		return
	}

	cardsView := components.ModalTaskView(task, taskTags, allTags)
	cardsView.Render(r.Context(), w)
}

// TODO: Optimize. Instead of reloading the whole table, update only a single row
func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, tags := resolveTaskFromForm(r)
	err := services.UpdateTask(task, tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	drawTaskTable(w, r)
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, tags := resolveTaskFromForm(r)
	err := services.SaveNewTask(task, tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	drawTaskTable(w, r)
}

func resolveTaskFromForm(r *http.Request) (models.Task, []models.TaskTag) {
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

	formCost := r.FormValue(consts.MODAL_TASK_COST_NAME)
	cost, err := models.StrToEnum[models.TaskCost](formCost)
	if err != nil {
		log.Printf("failed to parse: %v: %v", formCost, err)
	}

	// Parse checkbox values - they will be "on" if checked, or empty if unchecked
	wipValue := r.FormValue("task-wip") == "on"
	plannedValue := r.FormValue("task-planned") == "on"
	var completed time.Time
	if r.FormValue("task-completed") == "on" {
		completed = time.Now()
	} else {
		completed = models.NOT_COMPLETED
	}

	var taskTags []models.TaskTag
	for key, _ := range r.Form {
		tag, found := strings.CutPrefix(key, "tag-")
		if found {
			taskTags = append(taskTags, models.TaskTag(tag))
		}
	}

	return models.Task{
		Id:        r.FormValue("card-id"),
		Content:   r.FormValue("card-text"),
		Title:     r.FormValue("card-title"),
		Priority:  prio,
		Wip:       wipValue,
		Planned:   plannedValue,
		Impact:    impact,
		Completed: completed,
		Cost:      cost,
	}, taskTags
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
	drawTaskView(w, r)
}

func drawTaskViewBody(w http.ResponseWriter, r *http.Request) {
	cards, err := findTasksOrWriteError(w)
	if err != nil {
		return
	}

	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}
	body := components.TasksViewBody(cards, settings)
	body.Render(r.Context(), w)
}

func drawTaskView(w http.ResponseWriter, r *http.Request) {
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
