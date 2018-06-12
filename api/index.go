package api

import (
	"net/http"

	faktory "github.com/contribsys/faktory/client"
	"github.com/ilikebits/go-core/api"
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

		job := faktory.NewJob("index.Package", req.ImportPath)
		err := queue.Push(job)
		if err != nil {
			return nil, api.ErrorInternal(err)
		}
		return api.OK(struct{}{}), nil
	})
}
