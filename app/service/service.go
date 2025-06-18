package service

import (
	"app/config"
	"app/models"
	"time"
)

type tokenStorage interface {
	LoadToken() (string, time.Time)
	SaveToken(t string, expiredAt time.Time)
	DeleteToken()
}

type Service struct {
	cfg          config.Service
	tokenStorage tokenStorage

	OnNeedLogin func()
	Role        string

	tokens        models.Tokens
	stopOnRefresh func()
}

func NewService(cfg config.Service, tokenStorage tokenStorage) (*Service, error) {
	refreshToken, expiredAt := tokenStorage.LoadToken()

	s := &Service{
		cfg:          cfg,
		tokenStorage: tokenStorage,
		tokens: models.Tokens{
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: expiredAt,
		},
	}

	return s, nil
}
