package handlers

import (
	"log"
	"net/http"

	"github.com/inaryzen/prio_cards/consts"
	"github.com/inaryzen/prio_cards/services"
)

func PostToggleCompletedFilter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	filter := r.Form.Get(consts.COMPLETED_FILTER_NAME)
	value := filter != ""
	err = services.SetCompletedFilter(value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}

	drawCardTable(w, r)
}
