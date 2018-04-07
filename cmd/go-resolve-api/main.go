package main

import (
	"log"
	"net/http"
	"os"

	"github.com/contribsys/faktory"
	"github.com/jmoiron/sqlx"

	"github.com/fossas/go-resolve/api"
)

func main() {
	log.Println("Starting up...")
	queue, err := faktory.Open()
	if err != nil {
		log.Fatalf("Could not connect to faktory: %s", err.Error())
	}
	defer queue.Close()
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err.Error())
	}
	defer db.Close()

	addr := ":80"
	mux := api.New(client)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Listening at %#v.\n", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err.Error())
	}
}
