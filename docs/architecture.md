ticket-booking-system/
├── cmd/
│   ├── api/                          # Entrypoint HTTP Server (REST API)
│   │   └── main.go
│   ├── worker/                       # Entrypoint Background Consumer (Kafka Processors)
│   │   └── main.go
│   └── migration/                    # Entrypoint Database Migration Runner
│       └── main.go
│
├── internal/                         # Logic bisnis privat
│   ├── middleware/                   # Global middlewares (Auth JWT, CORS, Rate Limit, Recovery)
│   │   ├── auth.go
│   │   ├── rate_limiter.go
│   │   └── recovery.go
│   │
│   ├── booking/                      # Domain: Process Checkout & Orders
│   │   ├── delivery/
│   │   │   ├── http/                 # Controller REST API
│   │   │   │   └── booking_handler.go
│   │   │   └── messaging/            # Kafka Consumer (Async DB Writes)
│   │   │       └── booking_consumer.go
│   │   ├── domain/                   # Entities, Value Objects, Business Rules
│   │   │   ├── booking.go
│   │   │   └── errors.go
│   │   ├── dto/                      # Struct Payload Request/Response
│   │   │   └── booking_dto.go
│   │   ├── queries/                  # <-- SQL Queries murni untuk sqlc
│   │   │   └── booking.sql
│   │   ├── repository/               # Database Access Interface & Implementation
│   │   │   ├── booking_repository.go
│   │   │   ├── postgres_booking_repo.go # Wrapper sqlc + pgxpool
│   │   │   └── db/                   # <-- Generated Code dari sqlc (pgx/v5)
│   │   │       ├── db.go
│   │   │       ├── models.go
│   │   │       └── booking.sql.go
│   │   └── service/                  # Business Use Cases (Publish Kafka Event)
│   │       └── booking_service.go
│   │
│   ├── seat/                         # Domain: Inventory & High Concurrency Locking
│   │   ├── delivery/
│   │   │   └── http/                 # Controller Hold Seat
│   │   │       └── seat_handler.go
│   │   ├── domain/
│   │   │   └── seat.go
│   │   ├── queries/                  # <-- SQL Queries murni untuk sqlc
│   │   │   └── seat.sql
│   │   ├── repository/
│   │   │   ├── seat_repository.go
│   │   │   ├── postgres_seat_repo.go # Implementasi pgxpool + sqlc
│   │   │   ├── redis_seat_repo.go    # Redis Lua Scripting untuk Atomic Seat Reserve
│   │   │   └── db/                   # <-- Generated Code dari sqlc
│   │   │       ├── db.go
│   │   │       ├── models.go
│   │   │       └── seat.sql.go
│   │   └── service/                  # Logic Hold/Release Seat via Distributed Lock
│   │       └── seat_service.go
│   │
│   ├── payment/                      # Domain: Payment Processing & Webhooks
│   │   ├── delivery/
│   │   │   ├── http/                 # Payment Webhook Handler (Gateway Callback)
│   │   │   └── messaging/            # Consumer Payment Timeout / Expiration
│   │   ├── domain/
│   │   ├── queries/
│   │   ├── repository/
│   │   └── service/
│   │
│   └── user/                         # Domain: User Authentication & Profile
│       ├── delivery/
│       │   └── http/
│       ├── domain/
│       ├── queries/
│       ├── repository/
│       └── service/
│
├── pkg/                              # Reusable library/utilitas teknis murni
│   ├── redislock/                    # Distributed Lock Manager (Mencegah Overselling)
│   ├── ratelimiter/                  # In-memory / Redis Rate Limiting
│   ├── logger/                       # Structured Logger (Zap / Zerolog)
│   ├── validator/                    # Request Payload Validator
│   └── response/                     # Standardized HTTP JSON Response Format
│
├── platform/                         # Technical Drivers & Low-level Infrastructure
│   ├── database/
│   │   ├── postgres.go              # Setup *pgxpool.Pool (pgx/v5)
│   │   └── redis.go                 # Setup Redis Client
│   ├── messaging/
│   │   ├── kafka_producer.go        # Low-level Kafka Event Publisher
│   │   └── kafka_consumer.go        # Low-level Kafka Group Consumer
│   └── config/
│       └── config.go                # Env/App Configuration Loader
│
├── api/                              # Kontrak/Skema API
│   ├── proto/                        # File .proto (gRPC)
│   └── openapi/                      # OpenAPI / Swagger Specification
│
├── migrations/                       # SQL Migration Scripts (Digunakan sqlc & db runner)
│   ├── 000001_create_seats_table.up.sql
│   └── 000002_create_bookings_table.up.sql
│
├── deployments/                      # Container & Infrastructure Orchestration
│   ├── docker/
│   │   ├── Dockerfile.api
│   │   └── Dockerfile.worker
│   ├── k8s/                          # Kubernetes Manifests
│   │   ├── deployment-api.yaml
│   │   ├── deployment-worker.yaml
│   │   └── hpa.yaml                  # Auto-scaling berdasarkan trafik
│   └── docker-compose.yaml           # Environment Dev (Postgres, Redis, Kafka, Kafdrop)
│
├── scripts/                          # Utility & Load Testing Scripts
│   └── loadtest/
│       └── war_ticket_scenario.js    # Script Load Test k6
│
├── test/                             # Integration & End-to-End Tests
│   ├── integration/
│   └── e2e/
│
├── sqlc.yaml                         # Konfigurasi utama sqlc generator (pgx/v5)
├── go.mod
├── go.sum
├── Makefile                          # Command Runner (make sqlc, make run-api, make migrate)
└── README.md