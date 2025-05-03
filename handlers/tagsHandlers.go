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

func DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tagName := r.PathValue("name")
	err := services.DeleteTag(models.TaskTag(tagName))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Re-render the tags list
	allTags, err := services.Tags()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Re-render the tags list content
	newContent := components.TagsListContent(models.EMPTY_TASK, nil, allTags)
	newContent.Render(r.Context(), w)
}
