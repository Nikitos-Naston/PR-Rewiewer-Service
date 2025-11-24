.PHONY: build

build:
	go build -v ./cmd/

up:
	migrate -path storage/migrations -database "postgres://localhost:5432/ReviewService?sslmode=disable&user=postgres&password=2909" up

down:
	migrate -path storage/migrations -database "postgres://localhost:5432/ReviewService?sslmode=disable&user=postgres&password=2909" down

run:
	./cmd.exe

.DEFAULT_GOAL := build