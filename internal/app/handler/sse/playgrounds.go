package sse

import (
	"io"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initPlaygrounds(sse *gin.RouterGroup) {
	go h.waitRefresh()

	sse.GET("/playgrounds", h.SSEServeMiddleware(), h.getPlaygrounds)
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

	clientChan, ok := value.(<-chan string)
	if !ok {
		return
	}

	c.SSEvent("get", h.service.Playgrounds.GetPlaygrounds())
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case _, ok := <-clientChan:
			if ok {
				c.SSEvent("get", h.service.Playgrounds.GetPlaygrounds())
			} else {
				return false
			}
		case <-c.Writer.CloseNotify():
			return false
		}

		return true
	})
}
