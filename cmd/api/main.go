package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/DimasFatchurroziq/ticket-booking-system/platform/config"
)

func main() {
	// 1. Inisialisasi Configuration & Core Logging
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	defer log.Sync() // Memastikan sisa buffer log terkirim saat aplikasi exit

	log.Info("Starting application setup...")

	// 2. Inisialisasi Infrastructure & Resource Connection
	db := config.NewDatabase(viperConfig, log)
	redisClient := config.NewRedisClient(viperConfig, log)
	producer := config.NewKafkaProducer(viperConfig, log)

	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	// 3. Wiring Dependency Injection & Routing melalui Bootstrap
	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validate:    validate,
		Config:      viperConfig,
		Producer:    producer,
		RedisClient: redisClient,
	})

	// 4. Running Fiber Server secara Asynchronous
	webPort := viperConfig.GetInt("WEB_PORT")
	if webPort == 0 {
		webPort = 8080 // Fallback default port jika key env kosong
	}

	go func() {
		log.Info("HTTP Server is listening 🚀", zap.Int("port", webPort))
		if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Fiber server crashed unexpectedly", zap.Error(err))
		}
	}()

	// 5. Graceful Shutdown Interceptor
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal // Menunggu sinyal OS masuk

	log.Info("Shutdown signal received. Starting graceful shutdown...")

	// Langkah A: Hentikan penerimaan HTTP Request baru & selesaikan request aktif
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Fiber forced to shutdown with error", zap.Error(err))
	} else {
		log.Info("Fiber server stopped gracefully")
	}

	// Langkah B: Tutup semua resource/koneksi infrastruktur SETELAH Fiber benar-benar berhenti
	log.Info("Closing background infrastructure connections...")

	if err := producer.Close(); err != nil {
		log.Error("Failed to close Kafka producer gracefully", zap.Error(err))
	}

	if err := redisClient.Close(); err != nil {
		log.Error("Failed to close Redis client gracefully", zap.Error(err))
	}

	db.Close() // Menutup connection pool Postgres

	log.Info("Server clean exit. Goodbye! 👋")
}
