package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sandbox-webservice/internal/data"
	"strconv"
	"time"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) getCreateEntityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		entities := []data.Entity{
			{
				ID:        1,
				CreatedAt: time.Now(),
				Name:      "name1",
				Labels:    []string{"Go lang", "Is", "all right"},
				Version:   1,
			},
			{
				ID:        2,
				CreatedAt: time.Now(),
				Name:      "name2",
				Labels:    []string{"Go lang2", "Is2", "all right2"},
				Version:   1,
			},
		}
		if err := app.writeJson(w, http.StatusOK, envelope{"entities": entities}); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		var input struct {
			ID     int64    `json:"id"`
			Name   string   `json:"name"`
			Labels []string `json:"labels"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "%v\n", input)
	}
}

func (app *application) getUpdateDeleteEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			app.getEntity(w, r)
		}
	case http.MethodPut:
		{
			app.updateEntity(w, r)
		}
	case http.MethodDelete:
		{
			app.deleteEntity(w, r)
		}
	default:
		{
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func (app *application) getEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/entities/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	entity := data.Entity{
		ID:        idInt,
		CreatedAt: time.Now(),
		Name:      "name",
		Labels:    []string{"Go lang", "Is", "all right"},
		Version:   1,
	}

	if err := app.writeJson(w, http.StatusOK, envelope{"entity": entity}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updateEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/entities/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	var input struct {
		Name   *string  `json:"name"`
		Labels []string `json:"labels"`
	}

	entity := data.Entity{
		ID:        idInt,
		CreatedAt: time.Now(),
		Name:      "name",
		Labels:    []string{"Go lang", "Is", "all right"},
		Version:   1,
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		entity.Name = *input.Name
	}

	if len(input.Labels) > 0 {
		entity.Labels = input.Labels
	}

	if err := app.writeJson(w, http.StatusOK, envelope{"entity": entity}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/entities/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Delete an entity with ID: %d", idInt)

}
