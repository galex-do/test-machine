# Test Management Platform - Database Migration Commands

.PHONY: help migrate migrate-up migrate-down migrate-status migrate-reset migrate-create

# Default environment variables
DATABASE_USER ?= postgres
DATABASE_PASSWORD ?= postgres
DATABASE_DB ?= test
DATABASE_HOST ?= localhost

# Database URL for local and container environments
DATABASE_URL_LOCAL = postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):5432/$(DATABASE_DB)?sslmode=disable
DATABASE_URL_CONTAINER = postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@postgres:5432/$(DATABASE_DB)?sslmode=disable

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Docker-based migration commands (recommended for consistency)
migrate: migrate-up ## Alias for migrate-up

migrate-up: ## Run all pending migrations using Docker
	@echo "Running database migrations using Docker..."
	docker-compose --profile migration up migrate --build

migrate-down: ## Rollback the last migration using Docker
	@echo "Rolling back last migration using Docker..."
	docker-compose run --rm migrate goose postgres "$(DATABASE_URL_CONTAINER)" down

migrate-status: ## Check migration status using Docker
	@echo "Checking migration status using Docker..."
	docker-compose run --rm migrate goose postgres "$(DATABASE_URL_CONTAINER)" status

migrate-reset: ## Reset database (rollback all migrations) using Docker
	@echo "WARNING: This will rollback ALL migrations!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose run --rm migrate goose postgres "$(DATABASE_URL_CONTAINER)" reset; \
	fi

migrate-create: ## Create a new migration file (Usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Please provide a migration name. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(NAME)"
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	filename="migrations/$${timestamp}_$(NAME).sql"; \
	echo "-- +goose Up" > $$filename; \
	echo "-- Add your up migration here" >> $$filename; \
	echo "" >> $$filename; \
	echo "-- +goose Down" >> $$filename; \
	echo "-- Add your down migration here" >> $$filename; \
	echo "Created migration file: $$filename"

# Local migration commands (requires local goose installation)
migrate-local-up: ## Run migrations using local goose installation
	@echo "Running migrations locally (requires goose to be installed)..."
	goose -dir migrations postgres "$(DATABASE_URL_LOCAL)" up

migrate-local-down: ## Rollback last migration using local goose installation
	@echo "Rolling back last migration locally..."
	goose -dir migrations postgres "$(DATABASE_URL_LOCAL)" down

migrate-local-status: ## Check migration status using local goose installation
	@echo "Checking migration status locally..."
	goose -dir migrations postgres "$(DATABASE_URL_LOCAL)" status

# Docker environment commands
up: ## Start all services (database + app)
	docker-compose up --build

up-db: ## Start only the database service
	docker-compose up postgres

down: ## Stop all services
	docker-compose down

clean: ## Stop services and remove volumes
	docker-compose down -v
	docker system prune -f

# Development helpers
logs: ## View application logs
	docker-compose logs -f app

logs-db: ## View database logs
	docker-compose logs -f postgres

shell-app: ## Connect to application container shell
	docker-compose exec app sh

shell-db: ## Connect to database container shell
	docker-compose exec postgres psql -U $(DATABASE_USER) -d $(DATABASE_DB)