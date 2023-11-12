package sse

import (
	"safechildhood/internal/app/service"
	"safechildhood/pkg/sse"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	server *sse.Server[string]

	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		server:  sse.NewServer[string](),
		service: service,
	}
}

func (h *Handler) Init(sse *gin.RouterGroup) {
	go h.server.Listen()

	h.initPlaygrounds(sse)
}
