package log

import (
	"context"
	"database/sql"
)

type DB struct {
	*sql.DB
}

func (ldb *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	From(ctx).Debug().Str("query", query).Msg("sending query")
	return ldb.DB.ExecContext(ctx, query, args...)
}

func (ldb *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	From(ctx).Debug().Str("query", query).Msg("sending query")
	return ldb.DB.QueryContext(ctx, query, args...)
}

func (ldb *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	From(ctx).Debug().Str("query", query).Msg("sending query")
	return ldb.DB.QueryRowContext(ctx, query, args...)
}
