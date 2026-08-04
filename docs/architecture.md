ticket-booking-system/
├── cmd/
│   ├── api/                        # Entrypoint HTTP Server (REST API)
│   │   └── main.go
│   ├── worker/                     # Entrypoint Background Consumer (Kafka Processors)
│   │   └── main.go
│   └── migration/                  # Entrypoint Database Migration Runner
│       └── main.go
│
├── internal/                       # Logic bisnis privat
│   ├── middleware/                 # Global middlewares (Auth JWT, CORS, Rate Limit, Recovery)
│   │   ├── auth.go
│   │   ├── rate_limiter.go
│   │   └── recovery.go
│   │
│   ├── booking/                    # Domain: Process Checkout & Orders
│   │   ├── delivery/
│   │   │   ├── http/               # Controller REST API & DTO
│   │   │   │   ├── booking_handler.go
│   │   │   │   └── booking_dto.go
│   │   │   └── messaging/          # Kafka Consumer (Async DB Writes)
│   │   │       └── booking_consumer.go
│   │   ├── domain/                 # BERSIH: Entities, Value Objects, Repository Interface
│   │   │   ├── booking.go
│   │   │   ├── repository.go       # Interface (Kontrak kosong)
│   │   │   └── errors.go
│   │   ├── service/                # Business Use Cases (Alur proses bisnis)
│   │   │   └── booking_service.go
│   │   └── repository/             # Implementasi Teknis Database & sqlc
│   │       ├── queries/            # <-- SQL Queries murni untuk sqlc dipindah ke sini
│   │       │   └── booking.sql
│   │       ├── db/                 # <-- Generated Code dari sqlc (pgx/v5)
│   │       │   ├── db.go
│   │       │   ├── models.go
│   │       │   └── booking.sql.go
│   │       ├── booking_repository.go
│   │       └── postgres_booking_repo.go # Mengimplementasikan interface domain pakai sqlc
│   │
│   ├── seat/                       # Domain: Inventory & High Concurrency Locking
│   │   ├── delivery/
│   │   │   └── http/
│   │   │       └── seat_handler.go
│   │   ├── domain/
│   │   │   ├── seat.go
│   │   │   └── repository.go
│   │   ├── service/
│   │   │   └── seat_service.go
│   │   └── repository/
│   │       ├── queries/            # <-- SQL Queries murni untuk sqlc
│   │       │   └── seat.sql
│   │       ├── db/                 # <-- Generated Code dari sqlc
│   │       │   ├── db.go
│   │       │   ├── models.go
│   │       │   └── seat.sql.go
│   │       ├── seat_repository.go
│   │       ├── postgres_seat_repo.go # Implementasi pgxpool + sqlc
│   │       └── redis_seat_repo.go   # Redis Lua Scripting untuk Atomic Seat Reserve
│   │
│   ├── payment/                    # Domain: Payment Processing & Webhooks
│   │   ├── delivery/
│   │   │   ├── http/
│   │   │   └── messaging/
│   │   ├── domain/
│   │   ├── service/
│   │   └── repository/
│   │       ├── queries/
│   │       ├── db/
│   │       └── ...
│   │
│   └── user/                       # Domain: User Authentication & Profile
│       ├── delivery/
│       │   └── http/
│       ├── domain/
│       ├── service/
│       └── repository/
│           ├── queries/
│           ├── db/
│           └── ...
│
├── pkg/                            # Reusable library/utilitas teknis murni
│   ├── redislock/                  # Distributed Lock Manager (Mencegah Overselling)
│   ├── ratelimiter/                # In-memory / Redis Rate Limiting
│   ├── logger/                     # Structured Logger (Zap / Zerolog)
│   ├── validator/                  # Request Payload Validator
│   └── response/                   # Standardized HTTP JSON Response Format
│
├── platform/                       # Technical Drivers & Low-level Infrastructure
│   ├── database/
│   │   ├── postgres.go             # Setup *pgxpool.Pool (pgx/v5)
│   │   └── redis.go                # Setup Redis Client
│   ├── messaging/
│   │   ├── kafka_producer.go       # Low-level Kafka Event Publisher
│   │   └── kafka_consumer.go       # Low-level Kafka Group Consumer
│   └── config/
│       └── config.go               # Env/App Configuration Loader
│
├── api/                            # Kontrak/Skema API
│   ├── proto/                      # File .proto (gRPC)
│   └── openapi/                    # OpenAPI / Swagger Specification
│
├── migrations/                     # SQL Migration Scripts (Digunakan sqlc & db runner)
│   ├── 000001_create_seats_table.up.sql
│   └── 000002_create_bookings_table.up.sql
│
├── deployments/                    # Container & Infrastructure Orchestration
│   ├── docker/
│   │   ├── Dockerfile.api
│   │   └── Dockerfile.worker
│   ├── k8s/                        # Kubernetes Manifests
│   │   ├── deployment-api.yaml
│   │   ├── deployment-worker.yaml
│   │   └── hpa.yaml                # Auto-scaling berdasarkan trafik
│   └── docker-compose.yaml         # Environment Dev (Postgres, Redis, Kafka, Kafdrop)
│
├── scripts/                        # Utility & Load Testing Scripts
│   └── loadtest/
│       └── war_ticket_scenario.js  # Script Load Test k6
│
├── test/                           # Integration & End-to-End Tests
│   ├── integration/
│   └── e2e/
│
├── sqlc.yaml                       # Konfigurasi utama sqlc generator (pgx/v5)
├── go.mod
├── go.sum
├── Makefile                        # Command Runner (make sqlc, make run-api, make migrate)
└── README.md







+--------+
| Client | (Mengirim HTTP Request: POST /api/v1/tickets + JSON Body & Header)
+--------+
    │
    ▼
+--------------------------------------------------------------------------+
| 1. DELIVERY LAYER (Fiber & Middleware)                                   |
|    - Router (route.go) menangkap URL endpoint                            |
|    - Auth Middleware memvalidasi Token (jika invalid -> 401 Unauthorized) |
|    - HTTP Handler menerima request, parsing ke Request DTO               |
+--------------------------------------------------------------------------+
    │
    ▼
+--------------------------------------------------------------------------+
| 2. USECASE / APPLICATION SERVICE LAYER                                   |
|    - Mengoordinasikan alur bisnis (Use Case)                             |
|    - Memanggil Domain Entity untuk validasi aturan bisnis                |
+--------------------------------------------------------------------------+
    │
    ▼
+--------------------------------------------------------------------------+
| 3. DOMAIN LAYER                                                          |
|    - Menjalankan logika murni bisnis (Rich Domain Model)                 |
|    - Memastikan aturan main terpenuhi (misal: kuota tiket, status valid) |
+--------------------------------------------------------------------------+
    │
    ▼
+--------------------------------------------------------------------------+
| 4. REPOSITORY & INFRASTRUCTURE LAYER                                     |
|    - Berkomunikasi dengan database PostgreSQL via pgxpool                |
|    - Menjalankan Query SQL (INSERT / SELECT)                             |
|    - (Opsional) Mengirim event ke Kafka Producer                         |
+--------------------------------------------------------------------------+
    │
    (Data berhasil disimpan / diambil)
    │
    ▼
+--------------------------------------------------------------------------+
| 5. KEMBALI KE USECASE & DELIVERY                                         |
|    - Repository mengembalikan data ke Usecase                            |
|    - Usecase merapikan data dan menyerahkannya ke HTTP Handler           |
|    - HTTP Handler memetakan data ke Response DTO                         |
+--------------------------------------------------------------------------+
    │
    ▼
+--------+
| Client | (Menerima HTTP Response: 201 Created + JSON Data)
+--------+