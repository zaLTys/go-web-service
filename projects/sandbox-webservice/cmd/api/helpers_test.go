package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadJSON(t *testing.T) {
	app := &application{}

	body := bytes.NewBufferString(`{"name": "Test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	rr := httptest.NewRecorder()

	var dst struct {
		Name string `json:"name"`
	}

	err := app.readJSON(rr, req, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst.Name != "Test" {
		t.Errorf("expected name 'Test', got %q", dst.Name)
	}
}

func TestWriteJSON(t *testing.T) {
	app := &application{}
	rr := httptest.NewRecorder()

	err := app.writeJSON(rr, http.StatusCreated, envelope{"msg": "ok"}, nil)
	if err != nil {
		t.Fatalf("unexpected error writing JSON: %v", err)
	}

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !contains(body, "ok") {
		t.Errorf("response body missing expected JSON element")
	}
}
