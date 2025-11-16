package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Sentinel errors
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Entity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Labels    []string  `json:"labels,omitempty"`
	Version   int32     `json:"-"`
	CreatedAt time.Time `json:"-"`
}

type EntityModel struct {
	DB *sql.DB
}

// -----------------------------------------------------------------------------
// Insert
// -----------------------------------------------------------------------------

func (m EntityModel) Insert(entity *Entity) error {
	query := `
		INSERT INTO entities (name, labels)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []any{
		entity.Name,
		pq.Array(entity.Labels),
	}

	return m.DB.QueryRow(query, args...).Scan(
		&entity.ID,
		&entity.CreatedAt,
		&entity.Version,
	)
}

// -----------------------------------------------------------------------------
// Get (single entity)
// -----------------------------------------------------------------------------

func (m EntityModel) Get(id int64) (*Entity, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, labels, version
		FROM entities
		WHERE id = $1`

	var entity Entity

	err := m.DB.QueryRow(query, id).Scan(
		&entity.ID,
		&entity.CreatedAt,
		&entity.Name,
		pq.Array(&entity.Labels),
		&entity.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &entity, nil
}

// -----------------------------------------------------------------------------
// Update (with optimistic concurrency)
// -----------------------------------------------------------------------------

func (m EntityModel) Update(entity *Entity) error {
	query := `
		UPDATE entities
		SET name = $1, labels = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []any{
		entity.Name,
		pq.Array(entity.Labels),
		entity.ID,
		entity.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&entity.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Delete
// -----------------------------------------------------------------------------

func (m EntityModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM entities
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// -----------------------------------------------------------------------------
// GetAll (basic version – no filters, pagination can be added later)
// -----------------------------------------------------------------------------

func (m EntityModel) GetAll() ([]*Entity, error) {
	query := `
		SELECT id, created_at, name, labels, version
		FROM entities
		ORDER BY id`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := []*Entity{}

	for rows.Next() {
		var entity Entity

		err := rows.Scan(
			&entity.ID,
			&entity.CreatedAt,
			&entity.Name,
			pq.Array(&entity.Labels),
			&entity.Version,
		)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
