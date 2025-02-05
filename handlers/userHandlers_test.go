package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/inaryzen/priotasks/consts"
)

func TestPostFilterName(t *testing.T) {
	tests := []struct {
		name           string
		filterName     string
		formData       string
		expectedStatus int
	}{
		{
			name:           "Set Completed Filter True",
			filterName:     consts.COMPLETED_FILTER_NAME,
			formData:       consts.COMPLETED_FILTER_NAME + "=true",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Set Completed From Date",
			filterName:     consts.FILTER_COMPLETED_FROM,
			formData:       consts.FILTER_COMPLETED_FROM + "=2024-01-01",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Set Completed To Date",
			filterName:     consts.FILTER_COMPLETED_TO,
			formData:       consts.FILTER_COMPLETED_TO + "=2024-02-01",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Date Format",
			filterName:     consts.FILTER_COMPLETED_FROM,
			formData:       consts.FILTER_COMPLETED_FROM + "=invalid-date",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unknown Filter",
			filterName:     "unknown-filter",
			formData:       "value=test",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request
			req := httptest.NewRequest(http.MethodPost, "/filter/"+tt.filterName, strings.NewReader(tt.formData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			// Execute request
			PostFilterName(rr, req)

			// Check response status
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
