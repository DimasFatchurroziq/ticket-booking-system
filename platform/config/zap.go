package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(viper *viper.Viper) *zap.Logger {
	var cfg zap.Config

	if viper.GetString("APP_ENV") == "development" {
		cfg = zap.NewDevelopmentConfig() // Menggunakan preset Dev (Text/Console, Human-readable)
	} else {
		cfg = zap.NewProductionConfig() // Menggunakan preset Prod (JSON format)
	}

	logLevel := zapcore.Level(viper.GetInt32("LOG_LEVEL"))
	cfg.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
