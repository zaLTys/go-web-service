package data

import (
	"database/sql"
	"testing"
)

func TestNewModels(t *testing.T) {
	db := &sql.DB{}
	models := NewModels(db)

	if models.Entities == nil {
		t.Fatal("expected Entities model to be initialized")
	}
}
