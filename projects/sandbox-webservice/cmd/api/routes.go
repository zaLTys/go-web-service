package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	// -------------------------------------------------------------------------
	// Middleware stack
	// -------------------------------------------------------------------------
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.enableCORS)          // (optional custom middleware)
	r.Use(app.rateLimitMiddleware) // (optional custom middleware)
	r.Use(app.authenticate)        // (optional auth)

	// -------------------------------------------------------------------------
	// Health Check
	// -------------------------------------------------------------------------
	r.Get("/v1/healthcheck", app.healthcheck)

	// -------------------------------------------------------------------------
	// Entity Routes (RESTful)
	// -------------------------------------------------------------------------
	r.Route("/v1/entities", func(r chi.Router) {

		r.Get("/", app.listEntities)
		r.Post("/", app.createEntity)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.getEntity)
			r.Put("/", app.updateEntity)
			r.Delete("/", app.deleteEntity)
		})
	})

	return r
}
