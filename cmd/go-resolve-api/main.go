package main

import (
	"flag"
	"net/http"
	"os"

	faktory "github.com/contribsys/faktory/client"
	"github.com/ilikebits/go-core/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fossas/go-resolve/api"
)

func main() {
	faktoryURL := flag.String("faktory", "", "faktory URL")
	pgURL := flag.String("db", "", "database URL")
	secret := flag.String("secret", "", "API secret")
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	log.Init(*debug)

	log.Debug().Msg("initializing server")
	log.Debug().Msg("connecting to faktory")
	os.Setenv("FAKTORY_URL", *faktoryURL)
	queue, err := faktory.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to faktory")
	}
	log.Debug().Msg("connecting to postgres")
	defer queue.Close()
	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to postgres")
	}
	defer db.Close()

	addr := ":80"
	mux := api.New(db, queue, *secret)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Debug().Msg("starting server")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("could not start server")
	}
}
