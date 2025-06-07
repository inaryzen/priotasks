package services

import (
	"strings"
	"testing"
	"time"

	"github.com/inaryzen/priotasks/models"
)

func Test_ExportTasksToYAML_Success(t *testing.T) {
	tasks := []models.Task{
		{
			Id:        "1",
			Title:     "Test Task",
			Content:   "Task content",
			Created:   time.Now(),
			Updated:   time.Now(),
			Completed: time.Time{},
			Priority:  models.PriorityHigh,
			Wip:       false,
			Planned:   true,
			Impact:    models.ImpactHigh,
			Cost:      models.CostS,
			Value:     0.8,
			Fun:       models.FunM,
		},
	}

	yamlData, err := ExportTasksToYAML(tasks)
	if err != nil {
		t.Errorf("ExportTasksToYAML failed: %v", err)
	}

	if len(yamlData) == 0 {
		t.Error("ExportTasksToYAML returned empty data")
	}

	// Convert to string for content checking
	yamlStr := string(yamlData)

	// Check for human-readable values
	expectedValues := []string{
		"priority: High", // PriorityHigh
		"impact: High",   // ImpactHigh
		"cost: S (~30m)", // CostS
		"fun: M",         // FunM
	}

	for _, expected := range expectedValues {
		if !strings.Contains(yamlStr, expected) {
			t.Errorf("YAML output missing expected value: %s", expected)
		}
	}
}

func Test_ExportTasksToYAML_TooManyTasks(t *testing.T) {
	tasks := make([]models.Task, 1001)
	_, err := ExportTasksToYAML(tasks)
	if err == nil {
		t.Error("ExportTasksToYAML should fail with too many tasks")
	}
}
