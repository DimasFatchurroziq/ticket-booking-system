-- ============================================
-- FILE: xxxxxx_your_migration_name.down.sql
-- ============================================

-- 10. DROP QUEUE_TOKENS
DROP INDEX IF EXISTS idx_queue_event_status;
DROP TABLE IF EXISTS queue_tokens;

-- 9. DROP BOOKING_STATUS_HISTORY
DROP INDEX IF EXISTS idx_status_history_booking;
DROP TABLE IF EXISTS booking_status_history;

-- 8. DROP PAYMENTS
DROP INDEX IF EXISTS idx_payments_booking;
DROP TABLE IF EXISTS payments;

-- 7. DROP BOOKING_SEATS
DROP INDEX IF EXISTS idx_booking_seats_booking;
DROP TABLE IF EXISTS booking_seats;

-- 6. DROP BOOKINGS
DROP INDEX IF EXISTS idx_bookings_status_expires;
DROP INDEX IF EXISTS idx_bookings_user;
DROP TABLE IF EXISTS bookings;

-- 5. DROP SEATS
DROP INDEX IF EXISTS idx_seats_held_until;
DROP INDEX IF EXISTS idx_seats_event_status;
DROP TABLE IF EXISTS seats;

-- 4. DROP SEAT_CATEGORIES
DROP TABLE IF EXISTS seat_categories;

-- 3. DROP EVENTS
DROP INDEX IF EXISTS idx_events_status_sale_start;
DROP TABLE IF EXISTS events;

-- 2. DROP VENUES
DROP TABLE IF EXISTS venues;

-- 1. DROP USERS
DROP TABLE IF EXISTS users;