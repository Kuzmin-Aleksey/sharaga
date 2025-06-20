package httpserver

import (
	"context"
	"net/http"
	"sharaga/internal/domain/entity"
)

type authService interface {
	Login(ctx context.Context, username, password string) (*entity.Tokens, error)
	UpdateTokens(ctx context.Context, refresh string) (*entity.Tokens, error)
}

type AuthServer struct {
	authService authService
}

func NewAuthServer(authService authService) *AuthServer {
	return &AuthServer{
		authService: authService,
	}
}

func (s *AuthServer) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	email := r.FormValue("email")
	password := r.FormValue("password")

	tokens, err := s.authService.Login(ctx, email, password)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, tokens, http.StatusOK)
}

func (s *AuthServer) refreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	refreshToken := r.FormValue("refresh_token")

	tokens, err := s.authService.UpdateTokens(ctx, refreshToken)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, tokens, http.StatusOK)
}
