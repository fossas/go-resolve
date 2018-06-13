package api

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/ilikebits/go-core/api"
	"github.com/ilikebits/go-core/log"
	"github.com/jmoiron/sqlx"

	"github.com/fossas/go-resolve/models"
)

type ResolveRequest struct {
	Packages []models.Package
	Secret   string
}

func HandleResolve(db *sqlx.DB, secret string) http.HandlerFunc {
	return api.Handle(func(r *api.Request) (*api.Response, *api.Error) {
		var req ResolveRequest
		apiErr := r.JSON(&req)
		if apiErr != nil {
			return nil, apiErr
		}
		logger := log.From(r.Context())
		arr := zerolog.Arr()
		for _, pkg := range req.Packages {
			arr = arr.Str(pkg.String())
		}
		logger.Debug().
			Str("Secret", req.Secret).
			Array("Packages", arr).
			Msg("HandleResolve")

		if req.Secret != secret {
			logger.Error().Msg("bad secret")
			return nil, &api.Error{
				HTTPStatusCode: http.StatusForbidden,
				ErrorCode:      "INCORRECT_SECRET",
				Message:        "incorrect API secret",
			}
		}

		// TODO: this is probably a performance bottleneck.
		for _, pkg := range req.Packages {
			_, err := db.ExecContext(r.Context(), `
				INSERT INTO packages (
					import_path, revision, hash, last_updated
				) VALUES (
					$1, $2, $3, now()
				) ON CONFLICT (import_path, revision) DO UPDATE SET
					hash = $3,
					last_updated = now()
			`, pkg.ImportPath, pkg.Revision, pkg.Hash)
			if err != nil {
				logger.Error().Err(err).Msg("insert error")
				return nil, api.ErrorInternal(err)
			}
		}
		return api.OK(nil), nil
	})
}
