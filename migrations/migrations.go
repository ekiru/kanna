package main

import (
	"log"

	"github.com/ekiru/kanna/db"
	"github.com/ekiru/kanna/db/migrations"
	"github.com/ekiru/kanna/models"
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
				Name:       "id",
				Type:       migrations.String,
				PrimaryKey: true,
				NotNull:    true,
			},
			migrations.Column{
				Name:    "name",
				Type:    migrations.String,
				NotNull: true,
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
		migrations.FreeForm{
			Identifier: "0002-create-example-actor",
			Upward: func(tx db.MigrationTx) {
				tx.Exec("insert into Actors (id, name, type, inbox, outbox) values (?, ?, ?, ?, ?)",
					"http://kanna.example/actor/srn", "srn", "Person",
					"http://kanna.example/actor/srn/inbox", "http://kanna.example/actor/srn/outbox",
				)
			},
			Downward: func(tx db.MigrationTx) {
				tx.Exec("delete Actors where name = ?", "srn")
			},
		},
		migrations.CreateTable(
			"0003-create-accounts-table",
			"Accounts",
			migrations.Column{
				Name:       "username",
				Type:       migrations.String,
				PrimaryKey: true,
				NotNull:    true,
			},
			migrations.Column{
				Name:    "passwordHash",
				Type:    migrations.String,
				NotNull: true,
			},
			migrations.Column{
				Name:    "passwordHashVersion",
				Type:    migrations.Int,
				NotNull: true,
			},
			migrations.Column{
				Name:    "actorId",
				Type:    migrations.String,
				NotNull: true,
				Unique:  true,
			},
		),
		migrations.FreeForm{
			Identifier: "0004-create-example-account",
			Upward: func(tx db.MigrationTx) {
				tx.Exec("insert into Accounts (username, passwordHash, passwordHashVersion, actorId) values (?, ?, ?, ?)",
					"srn", models.HashScrypt.Hash("examplePassword", nil), models.HashScrypt,
					"http://kanna.example/actor/srn",
				)
			},
			Downward: func(tx db.MigrationTx) {
				tx.Exec("delete Accounts where username = ?", "srn")
			},
		},
		migrations.CreateTable(
			"0005-create-posts-table",
			"Posts",
			migrations.Column{
				Name:       "id",
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
				Name:    "audience",
				Type:    migrations.String,
				NotNull: true,
			},
			migrations.Column{
				Name:    "authorId",
				Type:    migrations.String,
				NotNull: true,
			},
			migrations.Column{
				Name:    "content",
				Type:    migrations.String,
				NotNull: true,
			},
			migrations.Column{
				Name:    "published",
				Type:    migrations.String,
				NotNull: true,
			},
		),
		migrations.FreeForm{
			Identifier: "0006-create-example-post",
			Upward: func(tx db.MigrationTx) {
				tx.Exec("insert into Posts (id, type, audience, authorId, content, published) values (?, ?, ?, ?, ?, ?)",
					"http://kanna.example/post/1", "Note", "https://www.w3.org/ns/activitystreams#Public",
					"http://kanna.example/actor/srn", "This is an example post!!!", "2017-10-21T15:49:45Z",
				)
			},
			Downward: func(tx db.MigrationTx) {
				tx.Exec("delete Posts where id = ?", "http://kanna.example/post/1")
			},
		},
	}
}
