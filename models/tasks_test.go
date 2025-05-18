package models

import (
	"testing"
	"time"
)

func Test_CalculateTotalTime(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  int
	}{
		{
			name:  "empty tasks",
			tasks: []Task{},
			want:  0,
		},
		{
			name: "single task",
			tasks: []Task{
				{Cost: CostM}, // 1h = 60 minutes
			},
			want: 60,
		},
		{
			name: "multiple tasks with different costs",
			tasks: []Task{
				{Cost: CostXS},  // 10 minutes
				{Cost: CostS},   // 30 minutes
				{Cost: CostM},   // 60 minutes
				{Cost: CostL},   // 120 minutes
				{Cost: CostXL},  // 240 minutes
				{Cost: CostXXL}, // 480 minutes
			},
			want: 10 + 30 + 60 + 120 + 240 + 480, // 940 minutes
		},
		{
			name: "multiple tasks with same cost",
			tasks: []Task{
				{Cost: CostM}, // 60 minutes
				{Cost: CostM}, // 60 minutes
				{Cost: CostM}, // 60 minutes
			},
			want: 180, // 3 hours = 180 minutes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateTotalTime(tt.tasks); got != tt.want {
				t.Errorf("CalculateTotalTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FormatTotalTime(t *testing.T) {
	tests := []struct {
		name    string
		minutes int
		want    string
	}{
		{
			name:    "zero minutes",
			minutes: 0,
			want:    "0m",
		},
		{
			name:    "minutes only",
			minutes: 45,
			want:    "45m",
		},
		{
			name:    "one hour exactly",
			minutes: 60,
			want:    "1h 0m",
		},
		{
			name:    "hours and minutes",
			minutes: 125,
			want:    "2h 5m",
		},
		{
			name:    "multiple hours",
			minutes: 180,
			want:    "3h 0m",
		},
		{
			name:    "large number of hours",
			minutes: 600,
			want:    "10h 0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTotalTime(tt.minutes); got != tt.want {
				t.Errorf("FormatTotalTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_titleFromContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "basic title",
			content: "This is a simple title",
			want:    "This is a simple title",
		},
		{
			name:    "title with newline",
			content: "First line\nSecond line",
			want:    "First line",
		},
		{
			name:    "title with windows newline",
			content: "First line\r\nSecond line",
			want:    "First line",
		},
		{
			name:    "long content truncation",
			content: "This is a very long title that should be truncated because it exceeds the maximum length limit",
			want:    "This is a very long title that should be truncated because it ex",
		},
		{
			name:    "empty content",
			content: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := titleFromContent(tt.content); got != tt.want {
				t.Errorf("titleFromContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	baseTime := time.Now()
	originalTask := Task{
		Id:        "test-id",
		Title:     "Original Title",
		Content:   "Original Content",
		Created:   baseTime,
		Updated:   baseTime,
		Completed: NOT_COMPLETED,
		Priority:  PriorityMedium,
		Wip:       false,
		Planned:   false,
		Impact:    ImpactModerate,
		Cost:      CostM,
	}

	tests := []struct {
		name     string
		current  Task
		changes  Task
		expected Task
	}{
		{
			name:    "basic update",
			current: originalTask,
			changes: Task{
				Id:        "test-id",
				Title:     "Original Title",
				Content:   "Changed Content",
				Created:   baseTime,
				Priority:  PriorityMedium,
				Impact:    ImpactModerate,
				Cost:      CostM,
				Completed: NOT_COMPLETED,
			},
			expected: Task{
				Id:        "test-id",
				Title:     "Original Title",
				Content:   "Changed Content",
				Created:   baseTime,
				Priority:  PriorityMedium,
				Impact:    ImpactModerate,
				Cost:      CostM,
				Completed: NOT_COMPLETED,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.current.Update(tt.changes)

			// Verify ID, Title, Content, and Completed status
			if result.Id != tt.expected.Id {
				t.Errorf("Update() Id = %v, want %v", result.Id, tt.expected.Id)
			}
			if result.Title != tt.expected.Title {
				t.Errorf("Update() Title = %v, want %v", result.Title, tt.expected.Title)
			}
			if result.Content != tt.expected.Content {
				t.Errorf("Update() Content = %v, want %v", result.Content, tt.expected.Content)
			}
			if !result.Completed.Equal(tt.expected.Completed) {
				t.Errorf("Update() Completed = %v, want %v", result.Completed, tt.expected.Completed)
			}

			// Verify Updated time is set
			if result.Updated.IsZero() {
				t.Error("Update() Updated time should be set")
			}

			// Verify Created time is preserved
			if !result.Created.Equal(tt.expected.Created) {
				t.Errorf("Update() Created = %v, want %v", result.Created, tt.expected.Created)
			}
		})
	}
}
