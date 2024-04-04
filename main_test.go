package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SanjaySinghRajpoot/YTapi/controller"
)

func TestGetVideosEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/videos", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetVideosHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
