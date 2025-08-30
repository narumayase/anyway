package handler_test

import (
	"anyway/internal/domain"
	httpHandler "anyway/internal/interfaces/http/handler"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsecase is a mock implementation of domain.Usecase
type MockUsecase struct {
	mock.Mock
}

// Send mocks the Send method of domain.Usecase
func (m *MockUsecase) Send(ctx context.Context, message domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// SetupRouter sets up a gin router for testing
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

// TestNewHandler tests the NewHandler constructor
func TestNewHandler(t *testing.T) {
	// Create a mock usecase
	mockUsecase := new(MockUsecase)

	// Create a new handler instance
	handler := httpHandler.NewHandler(mockUsecase)

	// Assert that the handler is not nil
	assert.NotNil(t, handler)
}

// TestSendSuccess tests the Send method when the usecase succeeds
func TestSendSuccess(t *testing.T) {
	// Create a mock usecase
	mockUsecase := new(MockUsecase)

	// Expect Send to be called and return nil (success)
	mockUsecase.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	// Create a new handler instance
	handler := httpHandler.NewHandler(mockUsecase)

	// Setup gin router and recorder
	router := SetupRouter()
	router.POST("/send", handler.Send)

	// Create a sample request body
	message := domain.Message{
		Key:     "test-key",
		Headers: map[string]string{"header1": "value1"},
		Content: []byte("test-content"),
	}
	jsonBody, _ := json.Marshal(message)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the expected methods were called on the mock
	mockUsecase.AssertExpectations(t)
}

// TestSendInvalidJSON tests the Send method with invalid JSON input
func TestSendInvalidJSON(t *testing.T) {
	// Create a mock usecase (it should not be called)
	mockUsecase := new(MockUsecase)

	// Create a new handler instance
	handler := httpHandler.NewHandler(mockUsecase)

	// Setup gin router and recorder
	router := SetupRouter()
	router.POST("/send", handler.Send)

	// Create an invalid request body
	invalidJson := []byte(`{"key": "test-key", "content": "invalid`)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(invalidJson))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert the error message in the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request format")

	// Assert that no methods were called on the mock usecase
	mockUsecase.AssertExpectations(t)
}

// TestSendUsecaseError tests the Send method when the usecase returns an error
func TestSendUsecaseError(t *testing.T) {
	// Create a mock usecase
	mockUsecase := new(MockUsecase)

	// Define an error to be returned by the usecase's Send method
	expectedErr := errors.New("failed to send message via usecase")

	// Expect Send to be called and return an error
	mockUsecase.On("Send", mock.Anything, mock.Anything).Return(expectedErr).Once()

	// Create a new handler instance
	handler := httpHandler.NewHandler(mockUsecase)

	// Setup gin router and recorder
	router := SetupRouter()
	router.POST("/send", handler.Send)

	// Create a sample request body
	message := domain.Message{
		Key:     "test-key",
		Headers: map[string]string{"header1": "value1"},
		Content: []byte("test-content"),
	}
	jsonBody, _ := json.Marshal(message)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the error message in the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Error processing message")

	// Assert that the expected methods were called on the mock
	mockUsecase.AssertExpectations(t)
}
