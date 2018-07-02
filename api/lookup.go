package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/ilikebits/go-core/api"
	"github.com/ilikebits/go-core/log"
	"github.com/jmoiron/sqlx"

	"github.com/fossas/go-resolve/models"
)

type LookupRequest struct {
	Hashes []string
}

type LookupResponse struct {
	Hashes   map[string]string
	Packages map[string]Package
}

type Package struct {
	Found      bool
	Ambiguous  bool
	Repository string
	Revision   Revision
	Candidates []Revision
}

type Revision struct {
	Hash      string
	Timestamp time.Time
	Matches   []string
}

// Lookup does multi-hash lookups.
//
// This implementation does revision selection in-memory, although it's possible
// to do this in SQL (e.g. by selecting revisions by hash, grouping by
// repository, and then intersecting).
func Lookup(db *sqlx.DB) http.HandlerFunc {
	return api.Handle(func(r *api.Request) api.Renderable {
		hash := chi.URLParam(r.Raw, "hash")

		logger := log.From(r.Context())
		logger.Debug().Str("hash", hash).Msg("Lookup")
		var pkg models.Package
		err := db.GetContext(r.Context(), &pkg, "SELECT * FROM packages WHERE hash = $1", hash)
		if err != nil {
			logger.Error().Err(err).Str("hash", hash).Msg("query error")
			return &api.Error{
				Raw:        err,
				StatusCode: http.StatusNotFound,
				ErrorCode:  "HASH_NOT_FOUND",
				Message:    "could not find hash",
			}
		}
		logger.Debug().
			Str("ImportPath", pkg.ImportPath).
			Str("Revision", pkg.Revision).
			Str("Version", pkg.Version).
			Msg("lookup succeeded")
		return api.OK( /* LookupResponse{Package: pkg} */ nil)
	})
}
