package rest

import (
	"safechildhood/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type Handler struct {
	router  *gin.Engine
	service *service.Service

	maxPhotos int

	textSanitazer *bluemonday.Policy
}

func New(router *gin.Engine, service *service.Service, maxPhotos int) *Handler {
	return &Handler{
		router:        router,
		service:       service,
		maxPhotos:     maxPhotos,
		textSanitazer: bluemonday.StrictPolicy(),
	}
}

func (h *Handler) Init() {
	h.initComplaints()
}
