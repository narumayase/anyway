package server

import (
	"anyway/config"
	"errors"
	"log"
	"net/http"

	"anyway/internal/domain"
	httphandler "anyway/internal/interfaces/http"
)

func Run(cfg config.Config, usecase domain.Usecase) {
	// Configure router
	router := httphandler.SetupRouter(usecase)

	// Start server
	serverAddr := ":" + cfg.Port
	if err := router.Run(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
