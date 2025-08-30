package http_test

import (
	"anyway/internal/domain"
	httpRouter "anyway/internal/interfaces/http"
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

// TestSetupRouterHealthCheck tests the /health endpoint
func TestSetupRouterHealthCheck(t *testing.T) {
	// Create a mock usecase (not used for health check, but required by SetupRouter)
	mockUsecase := new(MockUsecase)

	// Setup the router
	gin.SetMode(gin.TestMode)
	router := httpRouter.SetupRouter(mockUsecase)

	// Create a new HTTP request to the health endpoint
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])
	assert.Equal(t, "anyway API is running", response["message"])
}

// TestSetupRouterSendSuccess tests the /api/v1/send endpoint with a successful send
func TestSetupRouterSendSuccess(t *testing.T) {
	// Create a mock usecase
	mockUsecase := new(MockUsecase)

	// Expect Send to be called and return nil (success)
	mockUsecase.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	// Setup the router
	gin.SetMode(gin.TestMode)
	router := httpRouter.SetupRouter(mockUsecase)

	// Create a sample request body
	message := domain.Message{
		Key:     "test-key",
		Headers: map[string]string{"header1": "value1"},
		Content: []byte("test-content"),
	}
	jsonBody, _ := json.Marshal(message)

	// Create a new HTTP request to the send endpoint
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/send", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the expected methods were called on the mock
	mockUsecase.AssertExpectations(t)
}

// TestSetupRouterSendUsecaseError tests the /api/v1/send endpoint when usecase returns an error
func TestSetupRouterSendUsecaseError(t *testing.T) {
	// Create a mock usecase
	mockUsecase := new(MockUsecase)

	// Define an error to be returned by the usecase's Send method
	expectedErr := errors.New("failed to send message via usecase")

	// Expect Send to be called and return an error
	mockUsecase.On("Send", mock.Anything, mock.Anything).Return(expectedErr).Once()

	// Setup the router
	gin.SetMode(gin.TestMode)
	router := httpRouter.SetupRouter(mockUsecase)

	// Create a sample request body
	message := domain.Message{
		Key:     "test-key",
		Headers: map[string]string{"header1": "value1"},
		Content: []byte("test-content"),
	}
	jsonBody, _ := json.Marshal(message)

	// Create a new HTTP request to the send endpoint
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/send", bytes.NewBuffer(jsonBody))
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
