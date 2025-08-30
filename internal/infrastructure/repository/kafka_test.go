package repository_test

import (
	"anyway/internal/domain"
	"anyway/internal/infrastructure/repository"
	"context"
	"errors"
	"testing"

	kafka "github.com/narumayase/anysher/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnysherKafkaRepository is a mock implementation of the external kafka.Repository
// We define an interface here because we cannot mock a concrete type from an external module directly.
type MockAnysherKafkaRepository struct {
	mock.Mock
}

// Send mocks the Send method of kafka.Repository
func (m *MockAnysherKafkaRepository) Send(ctx context.Context, message kafka.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// Close mocks the Close method of kafka.Repository
func (m *MockAnysherKafkaRepository) Close() {
	m.Called()
}

// TestNewKafkaRepository tests the NewKafkaRepository constructor
func TestNewKafkaRepository(t *testing.T) {
	// Create a mock anysher kafka repository
	mockAnysherKafkaRepo := new(MockAnysherKafkaRepository)

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaRepo)

	// Assert that the repository is not nil
	assert.NotNil(t, kRepository)
}

// TestProduceSuccess tests the Produce method when anysher kafka Send succeeds
func TestProduceSuccess(t *testing.T) {
	// Create a mock anysher kafka repository
	mockAnysherKafkaRepo := new(MockAnysherKafkaRepository)

	// Expect Send to be called and return nil (success)
	mockAnysherKafkaRepo.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaRepo)

	// Define a sample domain message
	domainMessage := domain.Message{
		Key:     "test-key",
		Headers: map[string]string{"header1": "value1"},
		Content: []byte("test-content"),
	}

	// Call the Produce method
	err := kRepository.Produce(context.Background(), domainMessage)

	// Assert that no error is returned
	assert.NoError(t, err)

	// Assert that the expected methods were called on the mock
	mockAnysherKafkaRepo.AssertExpectations(t)
}

// TestProduceError tests the Produce method when anysher kafka Send returns an error
func TestProduceError(t *testing.T) {
	// Create a mock anysher kafka repository
	mockAnysherKafkaRepo := new(MockAnysherKafkaRepository)

	// Define an error to be returned by anysher kafka Send
	expectedErr := errors.New("failed to send message to anysher kafka")

	// Expect Send to be called and return an error
	mockAnysherKafkaRepo.On("Send", mock.Anything, mock.Anything).Return(expectedErr).Once()

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaRepo)

	// Define a sample domain message
	domainMessage := domain.Message{
		Key:     "test-key",
		Headers: map[string]string{"header1": "value1"},
		Content: []byte("test-content"),
	}

	// Call the Produce method
	err := kRepository.Produce(context.Background(), domainMessage)

	// Assert that the expected error is returned
	assert.EqualError(t, err, expectedErr.Error())

	// Assert that the expected methods were called on the mock
	mockAnysherKafkaRepo.AssertExpectations(t)
}

// TestClose tests the Close method
func TestClose(t *testing.T) {
	// Create a mock anysher kafka repository
	mockAnysherKafkaRepo := new(MockAnysherKafkaRepository)

	// Expect Close to be called
	mockAnysherKafkaRepo.On("Close").Return().Once()

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaRepo)

	// Call the Close method
	kRepository.Close()

	// Assert that the expected methods were called on the mock
	mockAnysherKafkaRepo.AssertExpectations(t)
}
