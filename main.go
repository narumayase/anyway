package main

import (
	"anyway/cmd/server"
	"anyway/config"
	"anyway/internal/application"
	"anyway/internal/infrastructure/repository"
	"github.com/narumayase/anysher/kafka"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	kafkaConfiguration := kafka.NewConfiguration(cfg.KafkaBroker, cfg.KafkaTopic, cfg.LogLevel)

	kafkaRepository, err := kafka.NewRepository(kafkaConfiguration)
	if err != nil {
		log.Fatal().Msgf("failed to create Kafka repository: %v", err)
	}

	// Create repository based on configuration
	producerRepository := repository.NewKafkaRepository(kafkaRepository)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create Kafka repository: %v", err)
	}

	// Create use case
	usecase := application.NewUsecase(producerRepository)

	server.Run(cfg, usecase)
}
