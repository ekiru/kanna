run:
	go run main.go

migrate:
	go run migrations/migrate.go

.PHONY: run migrate

