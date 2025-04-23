package server

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func (s *Server) ServerRun(handler http.Handler, port string) error {
	s.server = &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Print("Server run success from port:", port)
	return s.server.ListenAndServe()
}
