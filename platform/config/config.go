package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	RedisHost   string
	RedisPort   string
	KafkaHost   string
	KafkaPort   string
}

func LoadConfig() *Config {
	// Membaca file .env jika ada (jika di production/docker, env biasanya di-inject langsung)
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan Environment Variable sistem/docker")
	}

	// 1. App Port
	port := getEnv("APP_PORT", "8080")

	// 2. Database Configuration
	dbUser := getEnv("POSTGRES_USER", "ticket_user")
	dbPass := getEnv("POSTGRES_PASSWORD", "ticket_password")
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbName := getEnv("POSTGRES_DB", "ticket_booking_db")

	// Format DSN Postgres untuk Go (Driver pgx / lib/pq)
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	return &Config{
		AppPort:     port,
		DatabaseURL: dbURL,
		RedisHost:   getEnv("REDIS_HOST", "localhost"),
		RedisPort:   getEnv("REDIS_PORT", "6379"),
		KafkaHost:   getEnv("KAFKA_HOST", "localhost"),
		KafkaPort:   getEnv("KAFKA_PORT", "9092"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
