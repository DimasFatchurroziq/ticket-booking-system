package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewRedisClient(viper *viper.Viper, log *zap.Logger) *redis.Client {
	addr := viper.GetString("REDIS_ADDR")
	pass := viper.GetString("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	// Cek koneksi Redis saat startup
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Gagal terhubung ke Redis", zap.Error(err))
	}

	log.Info("Berhasil terhubung ke Redis")
	return client
}
