package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func DeleteTagName(w http.ResponseWriter, r *http.Request) {
	tagStr := r.PathValue("name")
	err := services.RemoveTagFromSettings(models.TaskTag(tagStr))
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
		return
	}

	r.ParseForm()

	filterName := r.PathValue("name")
	switch filterName {
	case consts.FILTER_NAME_HIDE_COMPLETED:
		filter := r.Form.Get(consts.FILTER_NAME_HIDE_COMPLETED)
		value := filter != ""
		t.FilterCompleted = value
	case consts.FILTER_NAME_HIDE_INCOMPLETED:
		filter := r.Form.Get(consts.FILTER_NAME_HIDE_INCOMPLETED)
		t.FilterIncompleted = filter != ""
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
	case consts.FILTER_NON_WIP:
		filter := r.Form.Get(consts.FILTER_NON_WIP)
		value := filter != ""
		t.FilterNonWip = value
	case consts.FILTER_PLANNED:
		filter := r.Form.Get(consts.FILTER_PLANNED)
		value := filter != ""
		t.Planned = value
	case consts.FILTER_NON_PLANNED:
		filter := r.Form.Get(consts.FILTER_NON_PLANNED)
		value := filter != ""
		t.NonPlanned = value
	case consts.FILTER_TAGS:
		tagStr := r.Form.Get(consts.FILTER_TAGS)
		if tagStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		t.Tags = append(t.Tags, models.TaskTag(tagStr))
	case consts.FILTER_SEARCH:
		searchText := r.Form.Get(consts.FILTER_SEARCH)
		t.SearchText = searchText
	case consts.FILTER_LIMIT_ENABLE:
		filter := r.Form.Get(consts.FILTER_LIMIT_ENABLE)
		value := filter != ""
		t.EnableLimit = value
	case consts.FILTER_LIMIT_COUNT:
		limitCountStr := r.Form.Get(consts.FILTER_LIMIT_COUNT)
		if limitCountStr != "" {
			limitCount, err := strconv.Atoi(limitCountStr)
			if err != nil || limitCount < 1 {
				postFilterNameError(w, filterName, fmt.Errorf("invalid limit count: %s", limitCountStr))
				return
			}
			t.LimitCount = limitCount
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("unknown filter name")
		return
	}

	s.TasksQuery = t
	err = services.UpdateUserSettings(s)
	if err != nil {
		internalServerError(w, err)
		return
	}

	common.Debug("PostFilterName: %v", s)
	common.Debug("PostFilterName: %v", t)
	common.Debug("PostFilterName: %v", filterName)

	w.WriteHeader(http.StatusOK)
	drawTaskViewBody(w, r)
}

func postFilterNameError(w http.ResponseWriter, filterName string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("PostFilterName: error updating filter %v: %v", filterName, err)
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("PostFilterName: internal error: %v", err)
}
