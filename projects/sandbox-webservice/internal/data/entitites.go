package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Entity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Labels    []string  `json:"labels,omitempty"`
	Version   int32     `json:"-"`
	CreatedAt time.Time `json:"-"` //wont be displayed after marshaling
}

type EntityModel struct {
	DB *sql.DB
}

func (b EntityModel) Insert(entity *Entity) error {
	query := `
		INSERT INTO entities (name, labels)
		VALUES ($1, $2)
		RETURNING id, created_at, version`

	args := []interface{}{entity.Name, pq.Array(entity.Labels)}
	// return the auto generated system values to Go object
	return b.DB.QueryRow(query, args...).Scan(&entity.ID, &entity.CreatedAt, &entity.Version)
}

func (b EntityModel) Get(id int64) (*Entity, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	query := `
		SELECT id, created_at, name, labels, version
		FROM entities
		WHERE id = $1`

	var entity Entity

	err := b.DB.QueryRow(query, id).Scan(
		&entity.ID,
		&entity.CreatedAt,
		&entity.Name,
		pq.Array(&entity.Labels),
		&entity.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}

	return &entity, nil
}

func (b EntityModel) Update(entity *Entity) error {
	query := `
		UPDATE entities
		SET name = $1, labels = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []interface{}{entity.Name, pq.Array(entity.Labels), entity.ID, entity.Version}
	return b.DB.QueryRow(query, args...).Scan(&entity.Version)
}

func (b EntityModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	query := `
		DELETE FROM entities
		WHERE id = $1`

	results, err := b.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (b EntityModel) GetAll() ([]*Entity, error) {
	query := `
	  SELECT * 
	  FROM entities
	  ORDER BY id`

	rows, err := b.DB.Query(query)
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
