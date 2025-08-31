package handler

import (
	"anyway/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Handler handles HTTP requests related to chat
type Handler struct {
	producerUsecase domain.Usecase
}

// NewHandler creates a new instance of the http handler
func NewHandler(usecase domain.Usecase) *Handler {
	return &Handler{
		producerUsecase: usecase,
	}
}

// Send processes the POST chat request
func (h *Handler) Send(c *gin.Context) {
	var request domain.Message

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}
	err := h.producerUsecase.Send(c.Request.Context(), request)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error processing message: " + err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
