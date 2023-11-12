package rest

import (
	"safechildhood/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

type Handler struct {
	service *service.Service

	maxPhotos int

	textSanitazer *bluemonday.Policy
}

func New(service *service.Service, maxPhotos int) *Handler {
	return &Handler{
		service:       service,
		maxPhotos:     maxPhotos,
		textSanitazer: bluemonday.StrictPolicy(),
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	h.initComplaints(api)
}
