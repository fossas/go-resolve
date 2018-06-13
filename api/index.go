package api

import (
	"net/http"

	faktory "github.com/contribsys/faktory/client"
	"github.com/ilikebits/go-core/api"
	"github.com/ilikebits/go-core/log"
)

type IndexRequest struct {
	ImportPath string
}

func HandleIndex(queue *faktory.Client) http.HandlerFunc {
	return api.Handle(func(r *api.Request) (*api.Response, *api.Error) {
		var req IndexRequest
		apiErr := r.JSON(&req)
		if apiErr != nil {
			return nil, apiErr
		}

		logger := log.From(r.Context())
		logger.Debug().Str("ImportPath", req.ImportPath).Msg("HandleIndex")
		job := faktory.NewJob("index.Package", req.ImportPath)
		logger.Debug().Str("JID", job.Jid).Msg("queuing index.Package job")
		err := queue.Push(job)
		if err != nil {
			logger.Error().Err(err).Str("JID", job.Jid).Msg("could not queue job")
			return nil, api.ErrorInternal(err)
		}
		return api.OK(nil), nil
	})
}
