package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	faktory "github.com/contribsys/faktory/client"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fossas/go-resolve/api"
)

func main() {
	faktoryURL := flag.String("faktory", "", "faktory URL")
	pgURL := flag.String("db", "", "database URL")
	secret := flag.String("secret", "", "API secret")
	flag.Parse()

	os.Setenv("FAKTORY_URL", *faktoryURL)
	queue, err := faktory.Open()
	if err != nil {
		log.Fatalf("Could not connect to faktory: %s", err.Error())
	}
	defer queue.Close()
	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err.Error())
	}
	defer db.Close()

	addr := ":80"
	mux := api.New(db, queue, *secret)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not initialize server: %s", err.Error())
	}
}
