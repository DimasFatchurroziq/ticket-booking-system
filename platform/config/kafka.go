package config

import (
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewKafkaConsumer(config *viper.Viper, log *zap.Logger) *kafka.Reader {
	brokers := strings.Split(config.GetString("KAFKA_BOOTSTRAP_SERVERS"), ",")

	readerConfig := kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  config.GetString("KAFKA_GROUP_ID"),
		Topic:    config.GetString("KAFKA_TOPIC"), // Segmentio menentukan topik di level Reader
		MinBytes: 10e3,                            // 10KB
		MaxBytes: 10e6,                            // 10MB
	}

	offsetReset := config.GetString("KAFKA_AUTO_OFFSET_RESET")
	if strings.ToLower(offsetReset) == "earliest" {
		readerConfig.StartOffset = kafka.FirstOffset
	} else {
		readerConfig.StartOffset = kafka.LastOffset
	}

	reader := kafka.NewReader(readerConfig)
	log.Info("Kafka Consumer (Reader) berhasil diinisialisasi")

	return reader
}

func NewKafkaProducer(config *viper.Viper, log *zap.Logger) *kafka.Writer {
	brokers := strings.Split(config.GetString("KAFKA_BOOTSTRAP_SERVERS"), ",")

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{}, // Algoritma distribusi pesan ke partisi
	}

	log.Info("Kafka Producer (Writer) berhasil diinisialisasi")

	return writer
}
