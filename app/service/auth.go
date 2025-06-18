package service

import (
	"app/models"
	"context"
	"time"
)

func (s *Service) Login(email, password string) error {
	tokens := new(models.Tokens)

	if err := s.request(tokens, post, urlLogin, nil,
		[2]string{"email", email},
		[2]string{"password", password},
	); err != nil {
		return err
	}

	s.tokens = *tokens

	s.tokenStorage.SaveToken(tokens.RefreshToken, tokens.RefreshTokenExpiresAt)

	ctx, cancel := context.WithCancel(context.Background())
	s.stopOnRefresh = cancel

	s.toRefresh(ctx)

	return nil
}

func (s *Service) RefreshTokens() {
	tokens := new(models.Tokens)

	if err := s.request(tokens, post, urlRefresh, nil,
		[2]string{"refresh_token", s.tokens.RefreshToken},
	); err != nil {
		s.OnNeedLogin()
		return
	}

	s.tokens = *tokens

	s.tokenStorage.SaveToken(tokens.RefreshToken, tokens.RefreshTokenExpiresAt)

	ctx, cancel := context.WithCancel(context.Background())
	s.stopOnRefresh = cancel

	s.toRefresh(ctx)
}

func (s *Service) Logout() error {
	s.tokenStorage.DeleteToken()
	return nil
}

func (s *Service) toRefresh(ctx context.Context) {
	if s.stopOnRefresh != nil {
		s.stopOnRefresh()
	}

	go func() {
		ticker := time.Tick(s.tokens.AccessTokenExpiresAt.Sub(time.Now()) - time.Minute)
		select {
		case <-ticker:
			s.RefreshTokens()
		case <-ctx.Done():
		}
	}()
}
