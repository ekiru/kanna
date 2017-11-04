run:
	go run main.go

generate:
	go generate github.com/ekiru/kanna/models

migrate:
	go run migrations/migrations.go

install-tools:
	go install github.com/ekiru/kanna/models/kanna-genmodel 

.PHONY: run generate migrate install-tools

