package sse

import "log"

type Server[T any] struct {
	message      chan T
	newClient    chan chan T
	closedClient chan chan T
	totalClients map[chan T]bool
}

func NewServer[T any]() *Server[T] {
	server := &Server[T]{
		message:      make(chan T),
		newClient:    make(chan chan T),
		closedClient: make(chan chan T),
		totalClients: make(map[chan T]bool),
	}

	return server
}

// This function returns new client and function that will close this client
//
//	clientChannel, closeClient := server.NewClient()
//	defer closeClient()
func (s *Server[T]) NewClient() (<-chan T, func()) {
	client := make(chan T)
	s.newClient <- client

	return client, func() {
		s.closedClient <- client
	}
}

func (s *Server[T]) WriteMessage(text T) {
	s.message <- text
}

func (s *Server[T]) Listen() {
	for {
		select {
		case client := <-s.newClient:
			log.Println("new")
			s.totalClients[client] = true

		case client := <-s.closedClient:
			log.Println("close")
			delete(s.totalClients, client)
			close(client)

		case message := <-s.message:
			for clientChan := range s.totalClients {
				clientChan <- message
			}
		}
	}
}
