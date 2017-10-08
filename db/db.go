// The db package handles connecting to Kanna's database and provides
// helper functions for database operations.
package db

import (
	"context"
	"database/sql"

	"github.com/ekiru/kanna/routes"
	_ "github.com/mattn/go-sqlite3"
)

// Open opens a connection to Kanna's database. This should only be
// called directly when running database migrations. In the normal
// operation of the application, the database will be passed through to
// request handlers via the context and can be accessed using the DB
// function.
func Open() (*sql.DB, error) {
	return sql.Open("sqlite3", "db.sqlite3")
}

// InitParams connects to the database and configures a Router to pass
// the database to request handlers via the context.
func InitParams(router *routes.Router) error {
	db, err := Open()
	if err != nil {
		return err
	}
	router.BaseParam(dbKey{}, db)
	return nil
}

type dbKey struct{}

// DB retrieves the database object from the request context.
func DB(ctx context.Context) *sql.DB {
	// TODO: maybe check this
	return ctx.Value(dbKey{}).(*sql.DB)
}
