package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fossas/go-resolve/api"
	"github.com/fossas/go-resolve/common"
	"github.com/go-redis/redis"
)

func main() {
	redisProvider := os.Getenv("REDIS_PROVIDER")
	redisURL := os.Getenv(redisProvider)
	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Could not initialize redis: %s", err.Error())
	}
	client := redis.NewClient(redisOptions)

	portProvider := os.Getenv("PORT_PROVIDER")
	port := os.Getenv(portProvider)
	_, err = strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Could not read port: %s", port)
	}
	addr := ":" + port
	mux := setupMux(client)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Listening at %s", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err.Error())
	}
}

func setupMux(cache *redis.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			res, err := json.Marshal(api.LookupResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				panic(err)
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}
		var req api.LookupRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			res, err := json.Marshal(api.LookupResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				panic(err)
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}

		pkg, err := lookup(cache, req.Hash)
		res, err := json.Marshal(api.LookupResponse{Ok: true, Result: pkg})
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
			res, err := json.Marshal(api.VerifyResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				panic(err)
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}
		var req api.VerifyRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			res, err := json.Marshal(api.VerifyResponse{Ok: false, Err: "INVALID_REQUEST"})
			if err != nil {
				panic(err)
			}
			http.Error(w, string(res), http.StatusBadRequest)
		}

		pkg, err := verify(cache, req.Package)
		res, err := json.Marshal(api.VerifyResponse{Ok: true, Expected: pkg})
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	return mux
}

var (
	errPackageNotFound       = errors.New("could not find package")
	errPackageChecksumFailed = errors.New("package checksum is different")
)

// Redis schema: string(Hash) => string(Name + " " + Revision)

func lookup(cache *redis.Client, hash string) (common.Package, error) {
	// TODO: implement this.
	return common.Package{}, nil
}

func verify(cache *redis.Client, actual common.Package) (common.Package, error) {
	// TODO: implement this.
	return common.Package{}, nil
}
