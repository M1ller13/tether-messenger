# Tether Messenger - Makefile for Docker operations

.PHONY: help build up down logs clean dev-up dev-down dev-logs test

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build all Docker images"
	@echo "  up        - Start all services in production mode"
	@echo "  down      - Stop all services"
	@echo "  logs      - Show logs from all services"
	@echo "  clean     - Remove all containers, volumes and images"
	@echo "  dev-up    - Start development environment (DB + Redis only)"
	@echo "  dev-down  - Stop development environment"
	@echo "  dev-logs  - Show development environment logs"
	@echo "  test      - Run tests"

# Production commands
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

# Development commands
dev-up:
	docker-compose -f docker-compose.dev.yml up -d

dev-down:
	docker-compose -f docker-compose.dev.yml down

dev-logs:
	docker-compose -f docker-compose.dev.yml logs -f

# Cleanup commands
clean:
	docker-compose down -v --rmi all
	docker-compose -f docker-compose.dev.yml down -v --rmi all
	docker system prune -f

# Test commands
test:
	cd server && go test ./...

# Database commands
db-migrate:
	cd server && go run main.go migrate

db-seed:
	cd server && go run main.go seed

# Development helpers
install-deps:
	cd client && npm install
	cd server && go mod download

# Quick start for development
dev-start: dev-up
	@echo "Development environment started!"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo ""
	@echo "Now you can run:"
	@echo "  Backend: cd server && go run main.go"
	@echo "  Frontend: cd client && npm run dev"
