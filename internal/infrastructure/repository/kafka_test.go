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

// MockAnysherKafkaClient is a mock implementation of the AnysherKafkaClient interface
type MockAnysherKafkaClient struct {
	mock.Mock
}

// Send mocks the Send method of AnysherKafkaClient
func (m *MockAnysherKafkaClient) Send(ctx context.Context, message kafka.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// Close mocks the Close method of AnysherKafkaClient
func (m *MockAnysherKafkaClient) Close() {
	m.Called()
}

// TestNewKafkaRepository tests the NewKafkaRepository constructor
func TestNewKafkaRepository(t *testing.T) {
	// Create a mock anysher kafka client
	mockAnysherKafkaClient := new(MockAnysherKafkaClient)

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaClient)

	// Assert that the repository is not nil
	assert.NotNil(t, kRepository)
}

// TestProduceSuccess tests the Produce method when anysher kafka Send succeeds
func TestProduceSuccess(t *testing.T) {
	// Create a mock anysher kafka client
	mockAnysherKafkaClient := new(MockAnysherKafkaClient)

	// Expect Send to be called and return nil (success)
	mockAnysherKafkaClient.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaClient)

	// Define a sample domain message
	domainMessage := domain.Message{
		Content: []byte("test-content"),
	}

	ctx := context.WithValue(context.Background(), "X-Correlation-Id", "test-correlation-id")
	ctx = context.WithValue(ctx, "X-Routing-Id", "test-routing-id")
	ctx = context.WithValue(ctx, "X-Request-Id", "test-request-id")

	// Call the Produce method
	err := kRepository.Produce(ctx, domainMessage)

	// Assert that no error is returned
	assert.NoError(t, err)

	// Assert that the expected methods were called on the mock
	mockAnysherKafkaClient.AssertExpectations(t)
}

// TestProduceError tests the Produce method when anysher kafka Send returns an error
func TestProduceError(t *testing.T) {
	// Create a mock anysher kafka client
	mockAnysherKafkaClient := new(MockAnysherKafkaClient)

	// Define an error to be returned by anysher kafka Send
	expectedErr := errors.New("failed to send message to anysher kafka")

	// Expect Send to be called and return an error
	mockAnysherKafkaClient.On("Send", mock.Anything, mock.Anything).Return(expectedErr).Once()

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaClient)

	// Define a sample domain message
	domainMessage := domain.Message{
		Content: []byte("test-content"),
	}

	// Call the Produce method
	err := kRepository.Produce(context.Background(), domainMessage)

	// Assert that the expected error is returned
	assert.EqualError(t, err, expectedErr.Error())

	// Assert that the expected methods were called on the mock
	mockAnysherKafkaClient.AssertExpectations(t)
}

// TestClose tests the Close method
func TestClose(t *testing.T) {
	// Create a mock anysher kafka client
	mockAnysherKafkaClient := new(MockAnysherKafkaClient)

	// Expect Close to be called
	mockAnysherKafkaClient.On("Close").Return().Once()

	// Create a new KafkaRepository instance
	kRepository := repository.NewKafkaRepository(mockAnysherKafkaClient)

	// Call the Close method
	kRepository.Close()

	// Assert that the expected methods were called on the mock
	mockAnysherKafkaClient.AssertExpectations(t)
}
