GOBIN=$(GOPATH)/bin
GOOSE=$(GOBIN)/goose
SQLBOILER=$(GOBIN)/sqlboiler
STRINGER=$(GOBIN)/stringer

MIGRATIONS_FOLDER=db/migrations
MIGRATIONS=$(wildcard $(MIGRATIONS_FOLDER)/*.sql)
DATABASE_URL=postgresql://local:hunter2@0.0.0.0:5432/goresolve?sslmode=disable
GO_FILES=$(shell find . -name "*.go")

all: $(GO_FILES) $(STRINGER)
	go generate ./...
	go install ./...

migrate-up: $(GOOSE) $(MIGRATIONS)
	cd $(MIGRATIONS_FOLDER); goose postgres "$(DATABASE_URL)" up

migrate-down: $(GOOSE) $(MIGRATIONS)
	cd $(MIGRATIONS_FOLDER); goose postgres "$(DATABASE_URL)" up

migrate-status: $(GOOSE) $(MIGRATIONS)
	cd $(MIGRATIONS_FOLDER); goose postgres "$(DATABASE_URL)" status

.PHONY: codegen
codegen: $(SQLBOILER) $(MIGRATIONS)
	sqlboiler postgres

.PHONY: dev
dev:
	sudo docker-compose up --force-recreate --build

.PHONY: bootstrap
bootstrap:
	sudo docker-compose up db
	make migrate-up
	make codegen

$(GOOSE): $(wildcard github.com/pressly/goose/**/*)
	go get -u github.com/pressly/goose/cmd/goose

$(SQLBOILER): $(wildcard github.com/volatiletech/sqlboiler/**/*)
	go get -u github.com/volatiletech/sqlboiler

$(STRINGER): $(wildcard golang.org/x/tools/cmd/stringer/**/*)
	go get -u golang.org/x/tools/cmd/stringer/**/*
