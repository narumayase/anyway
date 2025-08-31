package application

import (
	"anyway/internal/domain"
	"context"
	"github.com/rs/zerolog/log"
)

// UsecaseImpl implements Usecase
type UsecaseImpl struct {
	producerRepository domain.ProducerRepository
}

// NewUsecase creates a new instance of the usecase
func NewUsecase(producerRepository domain.ProducerRepository) domain.Usecase {
	return &UsecaseImpl{
		producerRepository: producerRepository,
	}
}

// Send sends the request
func (uc *UsecaseImpl) Send(ctx context.Context, message domain.Message) error {
	err := uc.producerRepository.Produce(ctx, message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return err
	}
	return nil
}
