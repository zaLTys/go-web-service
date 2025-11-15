package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func (app *application) getCreateEntityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Display a list of entities")
	}

	if r.Method == http.MethodPost {
		fmt.Fprintln(w, "added new entity")
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
	fmt.Fprintf(w, "Display the details of entity with ID: %d", idInt)
}

func (app *application) updateEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/entities/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Update the details of entity with ID: %d", idInt)

}

func (app *application) deleteEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/entities/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Delete an entity with ID: %d", idInt)

}
