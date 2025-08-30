package repository

import (
	"anyway/internal/domain"
	"context"
	"github.com/narumayase/anysher/kafka"
	"github.com/rs/zerolog/log"
)

// AnysherKafkaClient defines the methods used from the external kafka.Repository
type AnysherKafkaClient interface {
	Send(ctx context.Context, message kafka.Message) error
	Close()
}

// KafkaRepository implements the ProducerRepository interface for Kafka.
type KafkaRepository struct {
	kafkaClient AnysherKafkaClient
}

func NewKafkaRepository(kafkaClient AnysherKafkaClient) domain.ProducerRepository {
	return &KafkaRepository{
		kafkaClient: kafkaClient,
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
	if err := r.kafkaClient.Send(ctx, payload); err != nil {
		log.Err(err).Msg("Failed to send message to Kafka")
		return err
	}
	return nil
}

// Close closes the Kafka producer.
func (r *KafkaRepository) Close() {
	r.kafkaClient.Close()
}