package handlers

import (
	"log"
	"net/http"

	"github.com/inaryzen/priotasks/services"
)

func GetTasksYamlHandler(w http.ResponseWriter, r *http.Request) {
	settings, err := findSettingsOrWriteError(w)
	if err != nil {
		return
	}

	tasks, err := services.FindTasks(settings.TasksQuery)
	if err != nil {
		internalServerError(w, err)
		return
	}

	yamlData, err := services.ExportTasksToYAML(tasks)
	if err != nil {
		log.Printf("Failed to export tasks to YAML: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", "attachment; filename=\"tasks.yaml\"")
	_, err = w.Write(yamlData)
	if err != nil {
		log.Printf("Failed to write YAML response: %v", err)
	}
}
