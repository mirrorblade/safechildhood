package sse

import (
	"safechildhood/internal/app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	router *gin.Engine

	server *Server

	service *service.Service
}

func New(router *gin.Engine, service *service.Service) *Handler {
	return &Handler{
		router:  router,
		server:  NewServer(),
		service: service,
	}
}

func (h *Handler) Init() {
	go h.server.Listen()

	h.initPlaygroundsSSE()
}
