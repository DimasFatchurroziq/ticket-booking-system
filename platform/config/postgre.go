package config

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewDatabase(viper *viper.Viper, log *zap.Logger) *pgxpool.Pool {
	username := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetInt("POSTGRES_PORT")
	database := viper.GetString("POSTGRES_DB")
	// fmt.Println(password)

	idleConnection := viper.GetInt("POSTGRES_POOL_MIN")
	maxConnection := viper.GetInt("POSTGRES_POOL_MAX")
	maxLifeTimeConnection := viper.GetInt("POSTGRES_POOL_LIFETIME")

	// log.Info("Checking DB config",
	// 	zap.String("user", username),
	// 	zap.String("host", host),
	// 	zap.Int("port", port),
	// 	zap.String("db", database),
	// 	zap.Bool("has_password", password != ""),
	// )

	// Membentuk DSN secara aman menggunakan net/url
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
		Path:   database,
	}
	q := dbURL.Query()
	q.Set("sslmode", "disable")
	dbURL.RawQuery = q.Encode()

	dsn := dbURL.String()
	fmt.Println("DSN:", dsn)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("failed to parse postgresql config", zap.Error(err))
	}

	config.MaxConns = int32(maxConnection)
	config.MinConns = int32(idleConnection)
	config.MaxConnLifetime = time.Second * time.Duration(maxLifeTimeConnection)
	config.MaxConnIdleTime = 15 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("failed to connect to postgresql", zap.Error(err))
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal("failed to ping postgresql database", zap.Error(err))
	}

	log.Info("Successfully connected to PostgreSQL via pgxpool 🚀")

	return dbpool
}
