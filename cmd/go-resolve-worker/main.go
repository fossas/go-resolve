package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fossas/go-resolve/api"
	"github.com/fossas/go-resolve/index"
	"github.com/fossas/go-resolve/models"
)

func main() {
	faktoryURL := flag.String("faktory", "", "faktory URL")
	pgURL := flag.String("db", "", "database URL")
	apiURL := flag.String("api", "", "API URL")
	secret := flag.String("secret", "", "API secret")
	flag.Parse()

	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err.Error())
	}
	defer db.Close()

	os.Setenv("FAKTORY_URL", *faktoryURL)
	m := worker.NewManager()
	m.Concurrency = runtime.GOMAXPROCS(0)
	m.Queues = []string{"default"}

	u, err := url.Parse(*apiURL)
	if err != nil {
		log.Fatalf("Bad API URL: %s", err.Error())
	}
	endpoint, err := u.Parse("/api/resolve")
	if err != nil {
		log.Fatalf("Bad API endpoint: %s", err.Error())
	}

	// TODO: add logging, metrics, tracing, and health.
	m.Register("index.Package", func(ctx worker.Context, args ...interface{}) error {
		log.Printf("Starting job %s: resolve.Single(%#v, %#v)", ctx.Jid(), args[0])
		importpath := args[0].(string)

		err := index.Repository(importpath, func(pkgs []models.Package) error {
			body, err := json.Marshal(api.ResolveRequest{
				Packages: pkgs,
				Secret:   *secret,
			})
			if err != nil {
				return err
			}
			_, err = http.Post(endpoint.String(), "application/json", bytes.NewReader(body))
			return err
		})
		if err != nil {
			return err
		}

		return nil
	})

	log.Println("Ready.")
	m.Run()
}
