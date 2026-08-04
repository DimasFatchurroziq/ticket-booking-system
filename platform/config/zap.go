package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(viper *viper.Viper) *zap.Logger {

	var cfg zap.Config

	if viper.GetString("APP_ENV") == "development" {
		cfg = zap.NewDevelopmentConfig()

		// Stacktrace hanya muncul pada ERROR ke atas
		cfg.Development = false
		cfg.Encoding = "console"

	} else {
		cfg = zap.NewProductionConfig()
	}

	logLevel := zapcore.InfoLevel

	if viper.IsSet("LOG_LEVEL") {
		logLevel = zapcore.Level(viper.GetInt32("LOG_LEVEL"))
	}

	cfg.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := cfg.Build(
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	if err != nil {
		panic(err)
	}

	return logger
}
