package sse

import "github.com/gin-gonic/gin"

type Server struct {
	message      chan string
	newClient    chan chan string
	closedClient chan chan string
	totalClients map[chan string]bool
}

func NewServer() *Server {
	server := &Server{
		message:      make(chan string),
		newClient:    make(chan chan string),
		closedClient: make(chan chan string),
		totalClients: make(map[chan string]bool),
	}

	return server
}

func (s *Server) Listen() {
	for {
		select {
		case client := <-s.newClient:
			s.totalClients[client] = true

		case client := <-s.closedClient:
			delete(s.totalClients, client)
			close(client)

		case message := <-s.message:
			for clientChan := range s.totalClients {
				clientChan <- message
			}
		}
	}
}

func (s *Server) WriteMessage(text string) {
	s.message <- text
}

func (s *Server) ServeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChannel := make(chan string)
		s.newClient <- clientChannel

		defer func() {
			s.closedClient <- clientChannel
		}()

		c.Set("clientChannel", clientChannel)

		c.Next()
	}
}
