package main

import (
	"net/http"
	"net/http/httptest"
	"sandbox-webservice/internal/data"
	"testing"
)

type MockEntityRepository struct {
	InsertFunc func(*data.Entity) error
	GetFunc    func(int64) (*data.Entity, error)
	UpdateFunc func(*data.Entity) error
	DeleteFunc func(int64) error
	GetAllFunc func() ([]*data.Entity, error)
}

func (m MockEntityRepository) Insert(e *data.Entity) error        { return nil }
func (m MockEntityRepository) Get(id int64) (*data.Entity, error) { return nil, nil }
func (m MockEntityRepository) Update(e *data.Entity) error        { return nil }
func (m MockEntityRepository) Delete(id int64) error              { return nil }
func (m MockEntityRepository) GetAll() ([]*data.Entity, error)    { return m.GetAllFunc() }

func TestListEntities(t *testing.T) {
	mockRepo := MockEntityRepository{
		GetAllFunc: func() ([]*data.Entity, error) {
			return []*data.Entity{
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B"},
			}, nil
		},
	}

	app := &application{
		models: &data.Models{
			Entities: mockRepo,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/entities", nil)
	rr := httptest.NewRecorder()

	app.listEntities(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	body := rr.Body.String()
	if !contains(body, "A") || !contains(body, "B") {
		t.Errorf("response body missing entities")
	}
}
