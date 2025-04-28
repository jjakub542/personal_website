# Ładowanie zmiennych z pliku .env
ifneq ("$(wildcard .env)","")
	include .env
endif

# Zmienna z .env
DB_USER ?= jjakub542
DB_NAME_TEST ?= db_test

DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE=migrate -path migrations -database "$(DB_URL)"
TEST_DB_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TEST)?sslmode=disable
MIGRATE_TEST=migrate -path migrations -database "$(TEST_DB_URL)"

.PHONY: run test lint build migrate-up migrate-down migrate-force migrate-create reset-db db-console

# Build the application
build:
	@echo "Building..."
	
	
	@go build -o main cmd/server/main.go

# Run the application
run:
	@go run cmd/server/main.go


# Database create superuser
superuser:
	@go run cmd/superuser/main.go

# Test the application
test:
	docker compose exec postgres psql -U $(DB_USERNAME) -d postgres -c "CREATE DATABASE $(DB_NAME_TEST);"
	@echo "Migrating test DB..."
	@$(MIGRATE_TEST) drop -f
	@$(MIGRATE_TEST) up
	@echo "Running tests..."
	@go test ./tests/... -v || true
	docker compose exec postgres psql -U $(DB_USERNAME) -d postgres -c "DROP DATABASE $(DB_NAME_TEST);"


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

migrate-up:
	@echo "Running DB migrations..."
	@$(MIGRATE) up

migrate-down:
	@echo "Rolling back..."
	@$(MIGRATE) down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations $$name

reset-db:
	@echo "Resetting DB inside Docker..."
	docker exec -i $$(docker compose ps -q postgres) psql -U $(DB_USERNAME) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	docker exec -i $$(docker compose ps -q postgres) psql -U $(DB_USERNAME) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@$(MIGRATE) up

compose-up:
	docker compose up -d

compose-down:
	docker compose down

db-console:
	psql $(DB_URL)
