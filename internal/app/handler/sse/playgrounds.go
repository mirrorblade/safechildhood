package sse

import (
	"io"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initPlaygroundsSSE() {
	go h.waitRefresh()

	h.router.GET("/sse", SSEHeadersMiddleware(), h.server.ServeMiddleware(), h.getPlaygrounds)
}

func (h *Handler) waitRefresh() {
	for range h.service.Playgrounds.Refresh() {
		h.server.WriteMessage("refresh")
	}
}

func (h *Handler) getPlaygrounds(c *gin.Context) {
	value, exist := c.Get("clientChannel")
	if !exist {
		return
	}

	clientChan, ok := value.(chan string)
	if !ok {
		return
	}

	c.SSEvent("playgrounds", h.service.Playgrounds.GetPlaygrounds())

	c.Stream(func(w io.Writer) bool {
		if _, ok := <-clientChan; ok {
			c.SSEvent("playgrounds", h.service.Playgrounds.GetPlaygrounds())
		}

		return true
	})
}
