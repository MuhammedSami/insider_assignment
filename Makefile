include db.mk

APP_NAME := mailAPP

.PHONY: up down restart logs ps migrate build run test lint

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose down
	docker-compose up -d

ps:
	docker-compose ps


build:
	go build -o $(APP_NAME) ./...

run:
	go run ./...

test:
	@echo "Running tests..."
	go test ./...

lint:
	golangci-lint run ./...