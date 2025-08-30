package domain

import (
	"context"
)

// Usecase defines the interface for the use case
type Usecase interface {
	Send(ctx context.Context, message Message) error
}
