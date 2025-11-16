package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	app := &application{
		config: config{env: "test"},
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	rr := httptest.NewRecorder()

	app.healthcheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !contains(body, "available") {
		t.Errorf("response body missing expected content")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || (len(s) >= len(substr) && (func() bool {
		return string(s[:len(substr)]) == substr || contains(s[1:], substr)
	}())))
}
