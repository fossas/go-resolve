package api

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	faktory "github.com/contribsys/faktory/client"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// New sets up a new instance of http.ServeMux given a database and task queue.
func New(db *sqlx.DB, queue *faktory.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		var req LookupRequest
		ok := mustRead(w, r, &req)
		if !ok {
			return
		}

		var res LookupResponse
		found, err := lookupByHash(db, req.Hash)
		if err == nil {
			res = LookupResponse{Ok: true, Result: found}
		} else if err == sql.ErrNoRows {
			res = LookupResponse{Ok: false, Err: "NO_PACKAGE_FOUND"}
		} else {
			log.Panicf("Could not look up package by hash %#v: %s", req.Hash, err.Error())
		}

		mustReply(w, res)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		var req VerifyRequest
		ok := mustRead(w, r, &req)
		if !ok {
			return
		}

		var res VerifyResponse
		expected, err := lookupByKey(db, req.Package.Key)
		if err == nil {
			if req.Hash == expected.Hash {
				res = VerifyResponse{Ok: true, Expected: expected}
			} else {
				res = VerifyResponse{Ok: false, Expected: expected}
			}
		} else if errors.Cause(err) == sql.ErrNoRows {
			res = VerifyResponse{Ok: false, Err: "NO_PACKAGE_FOUND"}
		} else {
			log.Panicf("Could not look up package by hash %#v: %s", req.Hash, err.Error())
		}

		mustReply(w, res)
	})

	mux.HandleFunc("/crawl", handle(func(w http.ResponseWriter, r *http.Request) {
		var req CrawlRequest
		ok := mustRead(w, r, &req)
		if !ok {
			return
		}

		var res CrawlResponse
		job := faktory.NewJob("resolve.Single", req.Name, req.Revision)
		err := queue.Push(job)
		if err != nil {
			res = CrawlResponse{Ok: false, Err: err.Error()}
		} else {
			res = CrawlResponse{Ok: true}
		}

		mustReply(w, res)
	}))

	return mux
}

func handle(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(mustMarshal(ErrorResponse{Ok: false, Err: "INTERNAL_ERROR"}))
			}
		}()
		handler(w, r)
	}
}

func mustRead(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res := mustMarshal(ErrorResponse{Ok: false, Err: "INVALID_REQUEST"})
		http.Error(w, string(res), http.StatusBadRequest)
		return false
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("Failed to unmarshal request (%#v): %s", string(body), err.Error())
		res := mustMarshal(ErrorResponse{Ok: false, Err: "INVALID_REQUEST"})
		http.Error(w, string(res), http.StatusBadRequest)
		return false
	}

	return true
}

func mustReply(w http.ResponseWriter, res interface{}) {
	msg := mustMarshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(msg)
	if err != nil {
		log.Panicf("Could not write %#v: %s", string(msg), err.Error())
	}
}

func mustMarshal(data interface{}) []byte {
	msg, err := json.Marshal(data)
	if err != nil {
		log.Panicf("Could not marshal %#v: %s", string(msg), err.Error())
	}
	return msg
}
