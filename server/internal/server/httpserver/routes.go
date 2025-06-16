package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
	"sharaga/internal/domain/entity"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
	put  = http.MethodPut
	del  = http.MethodDelete
)

func (s *Server) RegisterRoutes(rtr *mux.Router) {
	rtr.HandleFunc("/login", s.authServer.login)

	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.NewUser, entity.UserRoleAdmin)).Methods(post)
	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.UpdateUser, entity.UserRoleAdmin)).Methods(put)
	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.GetAll, entity.UserRoleAdmin)).Methods(get)
	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.DeleteUser, entity.UserRoleAdmin)).Methods(del)

	rtr.HandleFunc("/order", s.authMw.withAuth(s.orderServer.NewOrder, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(post)
	rtr.HandleFunc("/order", s.authMw.withAuth(s.orderServer.GetAll, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)
	rtr.HandleFunc("/order/by-partner", s.authMw.withAuth(s.orderServer.GetByPartner, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)

}
