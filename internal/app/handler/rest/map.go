package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initMap() {
	interactive_map := h.router.Group("/map")
	{
		interactive_map.GET("/config", h.getConfig)
	}
}

func (h *Handler) getConfig(c *gin.Context) {
	c.JSON(http.StatusOK, h.handlerConfig.Map)
}
