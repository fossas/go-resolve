package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/url"
	"os"
	"runtime"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/ilikebits/go-core/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/fossas/go-resolve/api"
	"github.com/fossas/go-resolve/index"
	"github.com/fossas/go-resolve/models"
)

func main() {
	faktoryURL := flag.String("faktory", "", "faktory URL")
	pgURL := flag.String("db", "", "database URL")
	apiURL := flag.String("api", "", "API URL")
	secret := flag.String("secret", "", "API secret")
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

	u, err := url.Parse(*apiURL)
	if err != nil {
		log.Fatal().Err(err).Msg("bad API URL")
	}
	endpoint, err := u.Parse("/api/resolve")
	if err != nil {
		log.Fatal().Err(err).Msg("bad API endpoint")
	}

	// TODO: add logging, metrics, tracing, and health.
	log.Debug().Msg("registering index.Package task")
	m.Register("index.Package", func(ctx worker.Context, args ...interface{}) error {
		importpath := args[0].(string)

		logger := log.With().Str("JID", ctx.Jid()).Logger()
		logger.Debug().Str("ImportPath", importpath).Msg("starting index.Package")
		err := index.Repository(importpath, func(pkgs []models.Package) error {
			arr := zerolog.Arr()
			for _, pkg := range pkgs {
				arr = arr.Str(pkg.String())
			}
			logger.Debug().Array("Packages", arr).Str("Secret", *secret).Msg("starting revision package hash upload")

			body, err := json.Marshal(api.ResolveRequest{
				Packages: pkgs,
				Secret:   *secret,
			})
			if err != nil {
				logger.Error().Err(err).Msg("could not marshal ResolveRequest")
				return err
			}
			_, err = http.Post(endpoint.String(), "application/json", bytes.NewReader(body))
			if err != nil {
				logger.Error().Err(err).Msg("upload error")
				return err
			}
			return nil
		})
		if err != nil {
			logger.Error().Err(err).Msg("index.Package error")
			return err
		}

		return nil
	})

	log.Debug().Msg("running worker")
	m.Run()
}
