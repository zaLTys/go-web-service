package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sandbox-webservice/internal/data"

	"github.com/go-chi/chi/v5"
)

// -----------------------------------------------------------------------------
// Health Check
// -----------------------------------------------------------------------------

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// -----------------------------------------------------------------------------
// Entities: List
// -----------------------------------------------------------------------------

func (app *application) listEntities(w http.ResponseWriter, r *http.Request) {
	entities, err := app.models.Entities.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"entities": entities}, nil)
}

// -----------------------------------------------------------------------------
// Entities: Create
// -----------------------------------------------------------------------------

func (app *application) createEntity(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string   `json:"name"`
		Labels []string `json:"labels"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequest(w, err)
		return
	}

	entity := &data.Entity{
		Name:   input.Name,
		Labels: input.Labels,
	}

	if err := app.models.Entities.Insert(entity); err != nil {
		app.serverError(w, err)
		return
	}

	headers := http.Header{}
	headers.Set("Location", "/v1/entities/"+strconv.FormatInt(entity.ID, 10))

	app.writeJSON(w, http.StatusCreated, envelope{"entity": entity}, headers)
}

// -----------------------------------------------------------------------------
// Entities: Get
// -----------------------------------------------------------------------------

func (app *application) getEntity(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}

	entity, err := app.models.Entities.Get(id)
	if err == data.ErrRecordNotFound {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"entity": entity}, nil)
}

// -----------------------------------------------------------------------------
// Entities: Update
// -----------------------------------------------------------------------------

func (app *application) updateEntity(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}

	entity, err := app.models.Entities.Get(id)
	if err == data.ErrRecordNotFound {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	var input struct {
		Name   *string  `json:"name"`
		Labels []string `json:"labels"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequest(w, err)
		return
	}

	if input.Name != nil {
		entity.Name = *input.Name
	}

	if input.Labels != nil {
		entity.Labels = input.Labels
	}

	if err := app.models.Entities.Update(entity); err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"entity": entity}, nil)
}

// -----------------------------------------------------------------------------
// Entities: Delete
// -----------------------------------------------------------------------------

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		app.notFound(w)
		return
	}

	if err := app.models.Entities.Delete(id); err != nil {
		if err == data.ErrRecordNotFound {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "entity successfully deleted"}, nil)
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

// Extract URL param and parse it safely.
func getIDParam(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.ParseInt(idStr, 10, 64)
}
