package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fossas/go-resolve/resolve"
	"github.com/go-redis/redis"
)

func New(cache *redis.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			res, err := json.Marshal(LookupResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				panic(err)
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}
		var req LookupRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			res, err := json.Marshal(LookupResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				panic(err)
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}

		pkg, err := lookup(cache, req.Hash)
		res, err := json.Marshal(LookupResponse{Ok: true, Result: pkg})
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request: %s", err.Error())
			res, err := json.Marshal(VerifyResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				log.Panicf("Could not marshal error response: %s", err.Error())
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}
		var req VerifyRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("Failed to unmarshal request (%#v): %s", string(body), err.Error())
			res, err := json.Marshal(VerifyResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				log.Panicf("Could not marshal error response: %s", err.Error())
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}

		pkg, err := verify(cache, req.Package)
		res, err := json.Marshal(VerifyResponse{Ok: true, Expected: pkg})
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	mux.HandleFunc("/crawl", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request: %s", err.Error())
			res, err := json.Marshal(CrawlResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				log.Panicf("Could not marshal error response: %s", err.Error())
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}
		var req CrawlRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("Failed to unmarshal request (%#v): %s", string(body), err.Error())
			res, err := json.Marshal(CrawlResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				log.Panicf("Could not marshal error response: %s", err.Error())
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}

		_, err = resolve.Resolve(req.Name, req.Revision)
		res, err := json.Marshal(CrawlResponse{Ok: true})
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	return mux
}
