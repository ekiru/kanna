package db

import (
	"context"
	"database/sql"

	"github.com/ekiru/kanna/routes"
	_ "github.com/mattn/go-sqlite3"
)

func Open() (*sql.DB, error) {
	return sql.Open("sqlite3", "db.sqlite3")
}

func InitParams(router *routes.Router) error {
	db, err := Open()
	if err != nil {
		return err
	}
	router.BaseParam(dbKey{}, db)
	return nil
}

type dbKey struct{}

func DB(ctx context.Context) *sql.DB {
	// TODO: maybe check this
	return ctx.Value(dbKey{}).(*sql.DB)
}
