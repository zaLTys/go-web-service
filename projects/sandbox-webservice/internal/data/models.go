package data

import "database/sql"

type Models struct {
	Entities EntityRepository
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Entities: &EntityModel{DB: db},
	}
}
