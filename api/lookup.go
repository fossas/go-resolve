package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ilikebits/go-core/api"
	"github.com/ilikebits/go-core/log"
	"github.com/jmoiron/sqlx"

	"github.com/fossas/go-resolve/models"
)

type LookupResponse struct {
	Package models.Package
}

func HandleLookup(db *sqlx.DB) http.HandlerFunc {
	return api.Handle(func(r *api.Request) (*api.Response, *api.Error) {
		hash := chi.URLParam(r.Raw, "hash")

		logger := log.From(r.Context())
		logger.Debug().Str("hash", hash).Msg("HandleLookup")
		var pkg models.Package
		err := db.GetContext(r.Context(), &pkg, "SELECT * FROM packages WHERE hash = $1", hash)
		if err != nil {
			logger.Error().Err(err).Str("hash", hash).Msg("query error")
			return nil, &api.Error{
				Raw:            err,
				HTTPStatusCode: http.StatusNotFound,
				ErrorCode:      "HASH_NOT_FOUND",
				Message:        "could not find hash",
			}
		}
		logger.Debug().
			Str("ImportPath", pkg.ImportPath).
			Str("Revision", pkg.Revision).
			Str("Version", pkg.Version).
			Msg("lookup succeeded")
		return api.OK(LookupResponse{Package: pkg}), nil
	})
}
