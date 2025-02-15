package models

import "testing"

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
