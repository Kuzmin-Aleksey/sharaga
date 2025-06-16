package httpserver

import (
	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes(rtr *mux.Router) {
	rtr.HandleFunc("/login", s.authServer.login)
}
