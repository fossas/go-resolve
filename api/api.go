// Package api implements the API server.
package api

import (
	"net/http"
	"time"

	faktory "github.com/contribsys/faktory/client"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

// New sets up a new instance of chi.Mux given a database and task queue.
func New(db *sqlx.DB, queue *faktory.Client, secret string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // TODO: use our own logger that's structured
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/api/lookup/{hash}", HandleLookup(db))
	r.Post("/api/index", HandleIndex(queue))
	r.Post("/api/resolve", HandleResolve(db, secret))

	return r
}
