package services

import (
	"fmt"

	"github.com/inaryzen/priotasks/models"
	"gopkg.in/yaml.v3"
)

// ExportTasksToYAML exports the provided tasks to YAML format
func ExportTasksToYAML(tasks []models.Task) ([]byte, error) {
	if len(tasks) > 1000 {
		return nil, fmt.Errorf("maximum number of tasks (1000) exceeded: got %d", len(tasks))
	}

	data, err := yaml.Marshal(tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tasks to YAML: %w", err)
	}

	return data, nil
}
