// Package db implements application-level database abstractions. It provides
// implementations for different possible databases.
package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/fossas/go-resolve/log"
)

var (
	_ boil.Executor = &Logged{}
)

type Logged struct {
	*sqlx.DB
}

func New(db *sqlx.DB) Logged {
	return Logged{DB: db}
}

func (db *Logged) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	log.FromContext(ctx).WithFields(log.Fields{
		"query": query,
		"args":  args,
	}).Debug("sending query")
	return db.DB.ExecContext(ctx, query, args...)
}

func (db *Logged) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	log.FromContext(ctx).WithFields(log.Fields{
		"query": query,
		"args":  args,
	}).Debug("sending query")
	return db.DB.QueryContext(ctx, query, args...)
}

func (db *Logged) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	log.FromContext(ctx).WithFields(log.Fields{
		"query": query,
		"args":  args,
	}).Debug("sending query")
	return db.DB.QueryRowContext(ctx, query, args...)
}
