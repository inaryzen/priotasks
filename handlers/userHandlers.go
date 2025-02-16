package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/inaryzen/priotasks/common"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/models"
	"github.com/inaryzen/priotasks/services"
)

func PostPreparedQuery(w http.ResponseWriter, r *http.Request) {
	preparedQueryName := r.PathValue("name")
	err := services.ApplyPreparedQuery(preparedQueryName)
	if err != nil {
		internalServerError(w, err)
	}
	drawTaskViewBody(w, r)
}

func PostFilterName(w http.ResponseWriter, r *http.Request) {
	s, err := services.FindUserSettings()
	t := s.TasksQuery

	if err != nil {
		internalServerError(w, err)
	}

	r.ParseForm()

	filterName := r.PathValue("name")
	switch filterName {
	case consts.FILTER_NAME_HIDE_COMPLETED:
		filter := r.Form.Get(consts.FILTER_NAME_HIDE_COMPLETED)
		value := filter != ""
		t.FilterCompleted = value
		if err != nil {
			postFilterNameError(w, filterName, err)
			return
		}
	case consts.FILTER_NAME_HIDE_INCOMPLETED:
		filter := r.Form.Get(consts.FILTER_NAME_HIDE_INCOMPLETED)
		t.FilterIncompleted = filter != ""
		if err != nil {
			postFilterNameError(w, filterName, err)
			return
		}
	case consts.FILTER_COMPLETED_FROM:
		value := r.Form.Get(consts.FILTER_COMPLETED_FROM)
		var completedFrom time.Time = models.NOT_COMPLETED
		if value != "" {
			completedFrom, err = time.Parse(consts.DEFAULT_DATE_FORMAT, value)
			if err != nil {
				postFilterNameError(w, filterName, err)
				return
			}
		}
		t.CompletedFrom = completedFrom
	case consts.FILTER_COMPLETED_TO:
		value := r.Form.Get(consts.FILTER_COMPLETED_TO)
		var completedTo time.Time = models.NOT_COMPLETED
		if value != "" {
			completedTo, err = time.Parse(consts.DEFAULT_DATE_FORMAT, value)
			if err != nil {
				postFilterNameError(w, filterName, err)
				return
			}
		}
		t.CompletedTo = completedTo
	case consts.FILTER_WIP:
		filter := r.Form.Get(consts.FILTER_WIP)
		value := filter != ""
		t.FilterWip = value
		if err != nil {
			postFilterNameError(w, filterName, err)
			return
		}
	case consts.FILTER_NON_WIP:
		filter := r.Form.Get(consts.FILTER_NON_WIP)
		value := filter != ""
		t.FilterNonWip = value
		if err != nil {
			postFilterNameError(w, filterName, err)
			return
		}
	case consts.FILTER_PLANNED:
		filter := r.Form.Get(consts.FILTER_PLANNED)
		value := filter != ""
		t.Planned = value
		if err != nil {
			postFilterNameError(w, filterName, err)
			return
		}
	case consts.FILTER_NON_PLANNED:
		filter := r.Form.Get(consts.FILTER_NON_PLANNED)
		value := filter != ""
		t.NonPlanned = value
		if err != nil {
			postFilterNameError(w, filterName, err)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("unknown filter name")
	}

	s.TasksQuery = t
	services.UpdateUserSettings(s)

	common.Debug("PostFilterName: %v", s)
	common.Debug("PostFilterName: %v", t)
	common.Debug("PostFilterName: %v", filterName)

	drawTaskTable(w, r)
}

func postFilterNameError(w http.ResponseWriter, filterName string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("PostFilterName: error updating filter %v: %v", filterName, err)
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("PostFilterName: internal error: %v", err)
}
