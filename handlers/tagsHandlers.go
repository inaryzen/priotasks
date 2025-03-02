package handlers

import (
	"net/http"

	"github.com/inaryzen/priotasks/components"
	"github.com/inaryzen/priotasks/consts"
	"github.com/inaryzen/priotasks/models"
	"github.com/inaryzen/priotasks/services"
)

func PostTagsHandler(w http.ResponseWriter, r *http.Request) {
	formValue := r.FormValue(consts.INPUT_NAME_NEW_TAG)
	newTag := models.TaskTag(formValue)
	err := services.SaveTag(newTag)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	tagComponent := components.TaskModalTag(newTag, false)
	tagComponent.Render(r.Context(), w)
}
