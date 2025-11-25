ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DB_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

.PHONY: buildrun

build:
	go build -v ./cmd/

up:
	migrate -path storage/migrations -database $(DB_DSN) down up

down:
	migrate -path storage/migrations -database $(DB_DSN) down down

run:
	./cmd.exe

buildrun:
	go build -v ./cmd/
	./cmd.exe

fullrestart:
	migrate -path storage/migrations -database $(DB_DSN) down down
	migrate -path storage/migrations -database $(DB_DSN) down up
	go build -v ./cmd/
	./cmd.exe

start:
	migrate -path storage/migrations -database "postgres://localhost:5432/ReviewService?sslmode=disable&user=postgres&password=2909" up
	go build -v ./cmd/
	./cmd.exe

.DEFAULT_GOAL := build