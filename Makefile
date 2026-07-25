APP_NAME=ticket-booking-system

GO=go
DOCKER=docker
COMPOSE=docker compose


.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run          Run application"
	@echo "  make build        Build application"
	@echo "  make test         Run tests"
	@echo "  make fmt          Format Go code"
	@echo "  make lint         Run linter"
	@echo "  make clean        Clean build files"
	@echo "  make docker-up    Start docker services"
	@echo "  make docker-down  Stop docker services"


# =========================
# Development
# =========================

.PHONY: run
run:
	$(GO) run ./cmd/api


.PHONY: build
build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/api


.PHONY: test
test:
	$(GO) test -v ./...


.PHONY: fmt
fmt:
	$(GO) fmt ./...


.PHONY: tidy
tidy:
	$(GO) mod tidy



# =========================
# Quality
# =========================

.PHONY: lint
lint:
	golangci-lint run



# =========================
# Docker
# =========================

.PHONY: docker-up
docker-up:
	$(COMPOSE) up -d


.PHONY: docker-down
docker-down:
	$(COMPOSE) down


.PHONY: docker-restart
docker-restart:
	$(COMPOSE) restart



# =========================
# Database
# =========================

.PHONY: migrate-up
migrate-up:
	@echo "Run database migration up"


.PHONY: migrate-down
migrate-down:
	@echo "Run database migration down"



# =========================
# Cleanup
# =========================

.PHONY: clean
clean:
	rm -rf bin/


sqlc-user:
	sqlc generate -f sqlc/user.yaml

sqlc-seat:
	sqlc generate -f sqlc/seat.yaml

sqlc-booking:
	sqlc generate -f sqlc/booking.yaml

sqlc-payment:
	sqlc generate -f sqlc/payment.yaml

sqlc:
	make sqlc-user
	make sqlc-seat
	make sqlc-booking
	make sqlc-payment