-- ============================================
-- 1. USERS
-- ============================================
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    full_name       VARCHAR(255) NOT NULL,
    phone_number    VARCHAR(20) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================
-- 2. VENUES
-- ============================================
CREATE TABLE venues (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    address         TEXT,
    city            VARCHAR(100),
    total_capacity  INT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================
-- 3. EVENTS
-- ============================================
CREATE TABLE events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id        UUID NOT NULL REFERENCES venues(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    event_start     TIMESTAMPTZ NOT NULL,   -- simpan dalam UTC
    event_end       TIMESTAMPTZ NOT NULL,
    sale_start      TIMESTAMPTZ NOT NULL,   -- kapan penjualan tiket dibuka
    sale_end        TIMESTAMPTZ NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft'
                    CHECK (status IN ('draft', 'published', 'ongoing', 'completed', 'cancelled')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_status_sale_start ON events(status, sale_start);

-- ============================================
-- 4. SEAT_CATEGORIES (tier harga: VIP, Regular, Ekonomi, dll)
-- ============================================
CREATE TABLE seat_categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id),
    name            VARCHAR(100) NOT NULL,      -- 'VIP', 'Regular', dst
    price           NUMERIC(12,2) NOT NULL,
    total_quota     INT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================
-- 5. SEATS  <-- tabel paling krusial untuk concurrency
-- ============================================
CREATE TABLE seats (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id            UUID NOT NULL REFERENCES events(id),
    seat_category_id    UUID NOT NULL REFERENCES seat_categories(id),
    seat_code           VARCHAR(20) NOT NULL,   -- misal 'A1', 'B12'
    status              VARCHAR(20) NOT NULL DEFAULT 'available'
                        CHECK (status IN ('available', 'held', 'booked', 'blocked')),
    version             INT NOT NULL DEFAULT 0,  -- untuk optimistic locking
    held_by_booking_id  UUID,                    -- referensi booking yang sedang hold
    held_until          TIMESTAMPTZ,              -- TTL untuk seat reservation
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (event_id, seat_code)
);

-- Index paling penting: query "cari kursi available" akan sangat sering dipanggil
CREATE INDEX idx_seats_event_status ON seats(event_id, status);
-- Index untuk proses cleanup kursi yang held tapi TTL sudah habis
CREATE INDEX idx_seats_held_until ON seats(held_until) WHERE status = 'held';

-- ============================================
-- 6. BOOKINGS
-- ============================================
CREATE TABLE bookings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    event_id        UUID NOT NULL REFERENCES events(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                    CHECK (status IN (
                        'pending',            -- baru dibuat, kursi belum tentu ke-hold
                        'seat_held',          -- kursi sudah di-lock sementara
                        'awaiting_payment',
                        'confirmed',
                        'expired',            -- TTL habis sebelum bayar
                        'cancelled',
                        'refunded'
                    )),
    total_amount    NUMERIC(12,2) NOT NULL,
    idempotency_key VARCHAR(100) UNIQUE NOT NULL,  -- WAJIB: cegah duplicate booking dari retry
    expires_at      TIMESTAMPTZ,                    -- batas waktu bayar
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_status_expires ON bookings(status, expires_at)
    WHERE status IN ('pending', 'seat_held', 'awaiting_payment');

-- ============================================
-- 7. BOOKING_SEATS (relasi many-to-many booking <-> seat)
-- ============================================
CREATE TABLE booking_seats (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID NOT NULL REFERENCES bookings(id),
    seat_id         UUID NOT NULL REFERENCES seats(id),
    price_at_booking NUMERIC(12,2) NOT NULL,   -- snapshot harga saat itu (harga bisa berubah)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (seat_id)   -- satu kursi cuma boleh terikat 1 booking aktif
);

CREATE INDEX idx_booking_seats_booking ON booking_seats(booking_id);

-- ============================================
-- 8. PAYMENTS
-- ============================================
CREATE TABLE payments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id          UUID NOT NULL REFERENCES bookings(id),
    provider            VARCHAR(50) NOT NULL,       -- 'midtrans', 'xendit', dll
    provider_ref_id     VARCHAR(255) UNIQUE,        -- ID transaksi dari payment gateway
    idempotency_key     VARCHAR(100) UNIQUE NOT NULL,
    amount              NUMERIC(12,2) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'success', 'failed', 'refunded')),
    paid_at             TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_payments_booking ON payments(booking_id);

-- ============================================
-- 9. BOOKING_STATUS_HISTORY (audit trail — sangat membantu saat debugging)
-- ============================================
CREATE TABLE booking_status_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID NOT NULL REFERENCES bookings(id),
    from_status     VARCHAR(20),
    to_status       VARCHAR(20) NOT NULL,
    changed_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    note            TEXT
);

CREATE INDEX idx_status_history_booking ON booking_status_history(booking_id);

-- ============================================
-- 10. QUEUE_TOKENS (untuk virtual waiting room — opsional tahap awal)
-- ============================================
CREATE TABLE queue_tokens (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL REFERENCES events(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    token           VARCHAR(100) UNIQUE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'waiting'
                    CHECK (status IN ('waiting', 'active', 'expired')),
    issued_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    active_until    TIMESTAMPTZ
);

CREATE INDEX idx_queue_event_status ON queue_tokens(event_id, status, issued_at);