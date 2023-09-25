package rest

import (
	"safechildhood/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type Handler struct {
	router  *gin.Engine
	service *service.Service

	textSanitazer *bluemonday.Policy
}

func New(router *gin.Engine, service *service.Service) *Handler {
	return &Handler{
		router:        router,
		service:       service,
		textSanitazer: bluemonday.StrictPolicy(),
	}
}

func (h *Handler) Init() {
	h.initComplaints()
}
