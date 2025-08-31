package application_test

import (
	"anyway/internal/application"
	"anyway/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProducerRepository is a mock implementation of domain.ProducerRepository
type MockProducerRepository struct {
	mock.Mock
}

// Produce mocks the Produce method of ProducerRepository
func (m *MockProducerRepository) Produce(ctx context.Context, message domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// Close mocks the Close method of ProducerRepository
func (m *MockProducerRepository) Close() {
	m.Called()
}

// TestNewUsecase tests the NewUsecase constructor
func TestNewUsecase(t *testing.T) {
	// Create a mock producer repository
	mockRepo := new(MockProducerRepository)

	// Create a new use case instance
	usecase := application.NewUsecase(mockRepo)

	// Assert that the use case is not nil
	assert.NotNil(t, usecase)

	// Assert that the use case is of the expected type (UsecaseImpl)
	// This requires type assertion or reflection if UsecaseImpl is not exported.
	// For now, we'll just check if it's not nil.
}

// TestSendSuccess tests the Send method when Produce succeeds
func TestSendSuccess(t *testing.T) {
	// Create a mock producer repository
	mockRepo := new(MockProducerRepository)

	// Expect Produce to be called and return nil (success)
	mockRepo.On("Produce", mock.Anything, mock.Anything).Return(nil).Once()

	// Create a new use case instance
	usecase := application.NewUsecase(mockRepo)

	// Define a sample message
	message := domain.Message{
		Content: []byte("test-content"),
	}

	// Call the Send method
	err := usecase.Send(context.Background(), message)

	// Assert that no error is returned
	assert.NoError(t, err)

	// Assert that the expected methods were called on the mock
	mockRepo.AssertExpectations(t)
}

// TestSendProduceError tests the Send method when Produce returns an error
func TestSendProduceError(t *testing.T) {
	// Create a mock producer repository
	mockRepo := new(MockProducerRepository)

	// Define an error to be returned by Produce
	expectedErr := errors.New("failed to produce message")

	// Expect Produce to be called and return an error
	mockRepo.On("Produce", mock.Anything, mock.Anything).Return(expectedErr).Once()

	// Create a new use case instance
	usecase := application.NewUsecase(mockRepo)

	// Define a sample message
	message := domain.Message{
		Content: []byte("test-content")}

	// Call the Send method
	err := usecase.Send(context.Background(), message)

	// Assert that the expected error is returned
	assert.EqualError(t, err, expectedErr.Error())

	// Assert that the expected methods were called on the mock
	mockRepo.AssertExpectations(t)
}
