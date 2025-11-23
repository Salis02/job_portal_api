package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Router *chi.Mux
	Port   string
}

func NewServer(port string) *Server {
	r := chi.NewRouter()

	return &Server{
		Router: r,
		Port:   port,
	}
}

func (s *Server) Start() {
	log.Info().Msgf("Server running on port %s", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s.Router)
}
