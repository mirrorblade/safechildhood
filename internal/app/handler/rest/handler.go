package rest

import (
	"safechildhood/internal/app/config"
	"safechildhood/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type Handler struct {
	router  *gin.Engine
	service *service.Service

	handlerConfig config.HandlerConfig

	textSanitazer *bluemonday.Policy
}

func New(router *gin.Engine, service *service.Service, handlerConfig config.HandlerConfig) *Handler {
	return &Handler{
		router:        router,
		service:       service,
		handlerConfig: handlerConfig,
		textSanitazer: bluemonday.StrictPolicy(),
	}
}

func (h *Handler) Init() {
	h.initComplaints()
	h.initMap()
}
