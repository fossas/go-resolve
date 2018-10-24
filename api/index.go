package api

import (
	"net/http"

	"github.com/apex/log"
	faktory "github.com/contribsys/faktory/client"

	"github.com/fossas/go-resolve/api/internal/serve"
)

// IndexRequest contains the arguments for an indexing request.
type IndexRequest struct {
	ImportPath  string
	Revision    string
	IncludeDeps bool // Only valid when Revision is specified.
}

// Index queues indexing jobs for packages.
func Index(queue *faktory.Client) http.HandlerFunc {
	return api.Handle(func(r *api.Request) api.Renderable {
		var req IndexRequest
		apiErr := r.JSON(&req)
		if apiErr != nil {
			return apiErr
		}

		logger := log.From(r.Context())
		logger.Debug().Str("ImportPath", req.ImportPath).Msg("Index")
		job := faktory.NewJob("index.Package", req.ImportPath)
		logger.Debug().Str("JID", job.Jid).Msg("queuing index.Package job")
		err := queue.Push(job)
		if err != nil {
			logger.Error().Err(err).Str("JID", job.Jid).Msg("could not queue job")
			return api.ErrorInternal(err)
		}
		return api.OK(nil)
	})
}
