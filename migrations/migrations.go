package main

import (
	"log"

	"github.com/ekiru/kanna/db"
	"github.com/ekiru/kanna/db/migrations"
)

func main() {
	conn, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Exec(`create table if not exists Migrations (
	id text primary key not null
)`)
	if err != nil {
		log.Fatal(err)
	}
	for _, migration := range Migrations() {
		log.Printf("Applying %s\n", migration.ID())
		err := db.ApplyMigration(conn, migration)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Migrations() []db.Migration {
	return []db.Migration{
		migrations.CreateTable("0001-create-actors",
			"Actors",
			migrations.Column{
				Name:       "name",
				Type:       migrations.String,
				PrimaryKey: true,
				NotNull:    true,
			},
			migrations.Column{
				Name:    "type",
				Type:    migrations.String,
				NotNull: true,
			},
			migrations.Column{
				Name:    "inbox",
				Type:    migrations.String,
				NotNull: true,
			},
			migrations.Column{
				Name:    "outbox",
				Type:    migrations.String,
				NotNull: true,
			},
		),
	}
}
