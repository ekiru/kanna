run:
	go run main.go

migrate:
	go run migrations/migrations.go

.PHONY: run migrate

