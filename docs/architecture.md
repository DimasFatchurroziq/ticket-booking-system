ticket-booking-system/
├── cmd/
│   ├── api/                    # entrypoint HTTP server
│   │   └── main.go
│   ├── worker/                 # entrypoint consumer Kafka/RabbitMQ
│   │   └── main.go
│   └── migration/               # entrypoint db migration runner
│       └── main.go
│
├── internal/                    # kode privat, tidak bisa diimport project lain
│   ├── booking/                 # domain: booking
│   │   ├── domain/              # entity, value object, business rule murni
│   │   │   ├── booking.go
│   │   │   └── errors.go
│   │   ├── repository/          # interface + implementasi akses data
│   │   │   ├── booking_repository.go
│   │   │   └── postgres_booking_repo.go
│   │   ├── service/             # business logic / use case
│   │   │   └── booking_service.go
│   │   ├── handler/             # HTTP/gRPC handler (controller layer)
│   │   │   └── booking_handler.go
│   │   └── dto/                 # request/response struct
│   │       └── booking_dto.go
│   │
│   ├── event/                   # domain: event/konser/pertandingan
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── service/
│   │   └── handler/
│   │
│   ├── seat/                    # domain: kursi & inventory
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── service/
│   │   └── handler/
│   │
│   ├── payment/                 # domain: pembayaran
│   │   └── ...
│   │
│   └── user/                    # domain: user & auth
│       └── ...
│
├── pkg/                          # kode reusable, boleh diimport project lain
│   ├── redislock/                # implementasi distributed lock
│   ├── ratelimiter/
│   ├── logger/
│   ├── validator/
│   └── response/                 # standard API response wrapper
│
├── platform/                     # infrastruktur & konfigurasi teknis
│   ├── database/
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── messaging/
│   │   ├── kafka_producer.go
│   │   └── kafka_consumer.go
│   └── config/
│       └── config.go
│
├── api/                           # kontrak API
│   ├── proto/                     # file .proto untuk gRPC
│   └── openapi/                   # spesifikasi OpenAPI/Swagger
│
├── migrations/                    # SQL migration files
│   └── 000001_create_bookings_table.up.sql
│
├── deployments/                   # infra as code
│   ├── docker/
│   │   ├── Dockerfile.api
│   │   └── Dockerfile.worker
│   ├── k8s/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── hpa.yaml               # horizontal pod autoscaler
│   └── docker-compose.yaml        # untuk local dev
│
├── scripts/                        # helper scripts (seed data, load test, dll)
│   └── loadtest/
│       └── booking_scenario.js     # skenario k6
│
├── test/
│   ├── integration/
│   └── e2e/
│
├── docs/
│   └── architecture.md
│
├── go.mod
├── go.sum
├── Makefile
└── README.md