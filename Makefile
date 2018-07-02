MIGRATIONS_FOLDER=db/migrations
DATABASE_URL=postgresql://local:hunter2@0.0.0.0:5432/goresolve?sslmode=disable
GO_FILES=$(shell find . -name "*.go")

all: $(GO_FILES)
	go install ./...

migrate-up: $(wildcard $(MIGRATIONS_FOLDER)/*.sql)
	cd $(MIGRATIONS_FOLDER); goose postgres "$(DATABASE_URL)" up

migrate-down: $(wildcard $(MIGRATIONS_FOLDER)/*.sql)
	cd $(MIGRATIONS_FOLDER); goose postgres "$(DATABASE_URL)" up

migrate-status: $(wildcard $(MIGRATIONS_FOLDER)/*.sql)
	cd $(MIGRATIONS_FOLDER); goose postgres "$(DATABASE_URL)" status

codegen:
	sqlboiler postgres

dev:
	sudo docker-compose up --force-recreate --build
