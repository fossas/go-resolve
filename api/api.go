// Package api implements the API server.
package api

import (
	"net/http"

	faktory "github.com/contribsys/faktory/client"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// New sets up a new instance of chi.Mux given a database and task queue.
func New(db *sqlx.DB, queue *faktory.Client, secret string) http.Handler {
	r := chi.NewRouter()

	r.Get("/api/lookup/{hash}", HandleLookup(db))
	r.Post("/api/index", HandleIndex(queue))
	r.Post("/api/resolve", HandleResolve(db, secret))

	return r
}
