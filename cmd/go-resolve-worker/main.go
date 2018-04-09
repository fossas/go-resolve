package main

import (
	"log"
	"os"
	"runtime"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/fossas/go-resolve/resolve"
)

func main() {
	log.Println("Starting up...")
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err.Error())
	}
	defer db.Close()

	m := worker.NewManager()
	m.Concurrency = runtime.GOMAXPROCS(0)
	m.Queues = []string{"default"}

	m.Register("resolve.Single", func(ctx worker.Context, args ...interface{}) error {
		log.Printf("Starting job %s: resolve.Single(%#v, %#v)", ctx.Jid(), args[0], args[1])
		name := args[0].(string)
		revision := args[1].(string)

		// Resolve the package to a hash.
		pkg, err := resolve.Single(name, revision)
		if err != nil {
			return errors.Wrapf(err, "could not resolve single revision %s %s", name, revision)
		}

		// Add hash to database.
		_, err = db.Exec(`
			INSERT INTO revisions (
				package, revision, hash
			)	VALUES (
				$1, $2, $3
			) ON CONFLICT (package, revision) DO UPDATE SET
				hash = $3
		`, name, revision, pkg.Hash)
		if err != nil {
			return errors.Wrapf(err, "could not update database for jid %s", ctx.Jid())
		}

		return nil
	})

	log.Println("Ready.")
	m.Run()
}
