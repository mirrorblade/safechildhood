package sse

import (
	"safechildhood/internal/app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	router *gin.Engine

	service *service.Service
}

func New(router *gin.Engine, service *service.Service) *Handler {
	return &Handler{
		router:  router,
		service: service,
	}
}

func (h *Handler) Init() {
	h.initPlaygroundsSSE()
}
