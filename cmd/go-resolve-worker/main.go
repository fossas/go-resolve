// Package go-resolve-worker implements the indexing workers for go-resolve.
package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/apex/log"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/fossas/go-resolve/index"
)

func main() {
	faktoryURL := flag.String("faktory", "", "faktory URL")
	pgURL := flag.String("db", "", "database URL")
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	log.Init(*debug)

	log.Debug().Msg("initializing worker")
	log.Debug().Msg("connecting to postgres")
	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to postgres")
	}
	defer db.Close()

	os.Setenv("FAKTORY_URL", *faktoryURL)
	m := worker.NewManager()
	m.Concurrency = runtime.GOMAXPROCS(0)
	m.Queues = []string{"default"}

	// TODO: add logging, metrics, tracing, and health.
	log.Debug().Msg("registering index.Package task")
	m.Register("index.Package", func(ctx worker.Context, args ...interface{}) error {
		importpath := args[0].(string)

		logger := log.With().Str("JID", ctx.Jid()).Logger()
		logger.Debug().Str("ImportPath", importpath).Msg("starting index.Package")
		pkgs, err := index.Repository(importpath)
		if err != nil {
			logger.Error().Err(err).Msg("index.Repository error")
			return err
		}

		// TODO: re-implement this as a database update.
		// 	logger.Debug().Array("Packages", arr).Msg("starting revision package hash upload")
		_ = pkgs

		return nil
	})

	log.Debug().Msg("running worker")
	m.Run()
}
