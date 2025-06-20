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
	rtr.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rtr.HandleFunc("/auth/login", s.authServer.login)
	rtr.HandleFunc("/auth/refresh", s.authServer.refreshToken)

	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.New, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(post)
	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.Update, entity.UserRoleAdmin)).Methods(put)
	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.GetAll, entity.UserRoleAdmin)).Methods(get)
	rtr.HandleFunc("/users", s.authMw.withAuth(s.userServer.Delete, entity.UserRoleAdmin)).Methods(del)
	rtr.HandleFunc("/users/self", s.authMw.withAuth(s.userServer.GetSelf, entity.UserRoleAdmin, entity.UserRoleManager, entity.UserRoleWorker)).Methods(get)

	rtr.HandleFunc("/orders", s.authMw.withAuth(s.orderServer.New, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(post)
	rtr.HandleFunc("/orders", s.authMw.withAuth(s.orderServer.GetAll, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)
	rtr.HandleFunc("/orders/by-partner", s.authMw.withAuth(s.orderServer.GetByPartner, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)
	rtr.HandleFunc("/orders/discount", s.authMw.withAuth(s.orderServer.CalcDiscount, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)

	rtr.HandleFunc("/partners", s.authMw.withAuth(s.partnerServer.New, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(post)
	rtr.HandleFunc("/partners", s.authMw.withAuth(s.partnerServer.GetAll, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)
	rtr.HandleFunc("/partners", s.authMw.withAuth(s.partnerServer.Update, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(put)
	rtr.HandleFunc("/partners", s.authMw.withAuth(s.partnerServer.Delete, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(del)

	rtr.HandleFunc("/products", s.authMw.withAuth(s.productServer.New, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(post)
	rtr.HandleFunc("/products", s.authMw.withAuth(s.productServer.GetAll, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)
	rtr.HandleFunc("/products", s.authMw.withAuth(s.productServer.Update, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(put)
	rtr.HandleFunc("/products", s.authMw.withAuth(s.productServer.Delete, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(del)
	rtr.HandleFunc("/products/types", s.authMw.withAuth(s.productServer.NewType, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(post)
	rtr.HandleFunc("/products/types", s.authMw.withAuth(s.productServer.GetTypes, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(get)
	rtr.HandleFunc("/products/types", s.authMw.withAuth(s.productServer.UpdateType, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(put)
	rtr.HandleFunc("/products/types", s.authMw.withAuth(s.productServer.DeleteType, entity.UserRoleAdmin, entity.UserRoleManager)).Methods(del)
}
