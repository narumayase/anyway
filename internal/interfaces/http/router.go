package http

import (
	"anyway/internal/domain"
	"anyway/internal/interfaces/http/handler"
	"github.com/narumayase/anysher/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(chatUseCase domain.Usecase) *gin.Engine {
	router := gin.Default()

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.HeadersToContext())
	router.Use(middleware.RequestIDToLogger())

	// Create the controller
	chatHandler := handler.NewHandler(chatUseCase)

	// API routes group
	api := router.Group("/api/v1")
	api.POST("/send", chatHandler.Send)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "anyway API is running",
		})
	})
	return router
}
