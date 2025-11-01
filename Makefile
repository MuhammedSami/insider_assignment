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

run: up
	go run ./cmd/... --password=secret --interval=10s

test:
	@echo "Running tests..."
	go test ./...

lint:
	golangci-lint run ./...

swagger-ui:
	docker run --rm -p 8090:8080 -v `pwd`/docs:/usr/share/nginx/html/docs -e URL=/docs/swagger.yaml swaggerapi/swagger-ui