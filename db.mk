# we could use env vars but I will see if I have time I will fix static values here and read from env
DB_URL := postgres://appuser:secret@localhost:5432/appdb?sslmode=disable

migrate:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

migrate-version:
	migrate -path db/migrations -database "$(DB_URL)" version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=add_users_table"; \
	else \
		migrate create -ext sql -dir db/migrations -seq $(name); \
	fi