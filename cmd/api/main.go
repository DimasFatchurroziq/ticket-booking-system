package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/DimasFatchurroziq/ticket-booking-system/pkg/logger"
	"github.com/DimasFatchurroziq/ticket-booking-system/pkg/validator"
	"github.com/DimasFatchurroziq/ticket-booking-system/platform/config"
	"github.com/DimasFatchurroziq/ticket-booking-system/platform/database"
	"github.com/DimasFatchurroziq/ticket-booking-system/platform/messaging"

	// Import domain layer (Ganti nama module sesuai go.mod Anda)
	bookingHandler "github.com/DimasFatchurroziq/ticket-booking-system/internal/booking/handler"
	bookingRepo "github.com/DimasFatchurroziq/ticket-booking-system/internal/booking/repository"
	bookingService "github.com/DimasFatchurroziq/ticket-booking-system/internal/booking/service"

	eventHandler "github.com/DimasFatchurroziq/ticket-booking-system/internal/event/handler"
	eventRepo "github.com/DimasFatchurroziq/ticket-booking-system/internal/event/repository"
	eventService "github.com/DimasFatchurroziq/ticket-booking-system/internal/event/service"
)

func main() {
	// 1. GLOBAL INITIALIZATION & CONFIG
	cfg := config.LoadConfig()

	// Skala industri menggunakan structured logger (misal: Zap/Logrus) daripkg/
	appLogger := logger.NewLogger(cfg.AppEnv)

	// 2. PLATFORM INFRASTRUCTURE INITIALIZATION
	db, err := database.ConnectPostgres(cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to PostgreSQL: ", err)
	}
	defer db.Close()

	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis: ", err)
	}
	defer redisClient.Close()

	kafkaProducer, err := messaging.NewKafkaProducer(cfg.KafkaBrokers)
	if err != nil {
		appLogger.Fatal("Failed to init Kafka Producer: ", err)
	}
	defer kafkaProducer.Close()

	// 3. FIBER APP INITIALIZATION
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		// Menghubungkan custom global validator Anda dari pkg/
		StructValidator: validator.NewStructValidator(),
	})

	// 4. DEPENDENCY INJECTION & ROUTE REGISTRATION PER DOMAIN

	// --- Domain: Event ---
	evRepo := eventRepo.NewPostgresEventRepository(db)
	evService := eventService.NewEventService(evRepo, redisClient)
	eventHandler.RegisterRoutes(app, evService)

	// --- Domain: Booking ---
	bkRepo := bookingRepo.NewPostgresBookingRepository(db)
	bkService := bookingService.NewBookingService(bkRepo, evService, kafkaProducer)
	bookingHandler.RegisterRoutes(app, bkService)

	// 5. START SERVER (ASYNC)
	go func() {
		appLogger.Infof("Server %s is running on port %s", cfg.AppName, cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			appLogger.Fatalf("Server dynamic crash: %v", err)
		}
	}()

	// 6. OS GRACEFUL SHUTDOWN INTERCEPTOR
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	appLogger.Info("Shutdown signal received. Gracefully closing Fiber...")

	// Memberikan toleransi 10 detik agar proses transaksi tiket yang krusial selesai
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.Fatalf("Fiber forced to shutdown: %v", err)
	}

	appLogger.Info("Server clean exit. Goodbye.")
}
