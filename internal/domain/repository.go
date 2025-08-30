package domain

import (
	"context"
)

// ProducerRepository defines the interface for the producer repository for queue messages
type ProducerRepository interface {
	Produce(ctx context.Context, message Message) error
	Close()
}
