package repository

import (
	"anyway/internal/domain"
	"context"
	"github.com/narumayase/anysher/kafka"
	"github.com/rs/zerolog/log"
)

// KafkaRepository implements the ProducerRepository interface for Kafka.
type KafkaRepository struct {
	kafkaRepository *kafka.Repository
}

func NewKafkaRepository(kafkaRepository *kafka.Repository) domain.ProducerRepository {
	return &KafkaRepository{
		kafkaRepository: kafkaRepository,
	}
}

// Produce a message to a Kafka topic.
func (r *KafkaRepository) Produce(ctx context.Context, message domain.Message) error {
	// Create a payload
	payload := kafka.Message{
		Key:     message.Key,
		Headers: message.Headers,
		Content: message.Content,
	}
	// Send the message
	if err := r.kafkaRepository.Send(ctx, payload); err != nil {
		log.Err(err).Msg("Failed to send message to Kafka")
		return err
	}
	return nil
}

// Close closes the Kafka producer.
func (r *KafkaRepository) Close() {
	r.kafkaRepository.Close()
}
