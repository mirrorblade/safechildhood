package sse

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initPlaygroundsSSE() {
	h.router.GET("/sse", SSEHeadersMiddleware(), h.getPlaygrounds)
}

func (h *Handler) getPlaygrounds(c *gin.Context) {
	c.SSEvent("playgrounds", h.service.Playgrounds.GetPlaygrounds())

	c.Stream(func(w io.Writer) bool {
		if h.service.Playgrounds.CheckRefreshState() {
			c.SSEvent("playgrounds", h.service.Playgrounds.GetPlaygrounds())

			time.Sleep(1500 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
		}

		return true
	})
}
