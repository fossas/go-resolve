package api

import (
	"net/http"

	"github.com/ilikebits/go-core/api"
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

		if req.Secret != secret {
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
				) ON CONFLICT (package, revision) DO UPDATE SET
					hash = $3
					last_updated = now()
			`, pkg.ImportPath, pkg.Revision, pkg.Hash)
			if err != nil {
				return nil, api.ErrorInternal(err)
			}
		}
		return api.OK(struct{}{}), nil
	})
}
