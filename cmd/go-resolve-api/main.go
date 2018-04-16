package main

import (
	"net/http"
	"os"

	log "github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	faktory "github.com/contribsys/faktory/client"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fossas/go-resolve/api"
)

func main() {
	log.SetHandler(cli.Default)
	startupTrace := log.Trace("startup")
	faktoryTrace := log.Trace("connecting to faktory")
	queue, err := faktory.Open()
	if err != nil {
		log.Fatalf("Could not connect to faktory: %s", err.Error())
	}
	defer queue.Close()
	faktoryTrace.Stop(nil)
	pgTrace := log.Trace("connecting to postgres")
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err.Error())
	}
	defer db.Close()
	pgTrace.Stop(nil)

	addr := ":80"
	mux := api.New(db, queue)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	startupTrace.Stop(nil)
	log.WithField("address", addr).Info("listening")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err.Error())
	}
}
