package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"sharaga/internal/config"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/contextx"
	"sharaga/pkg/failure"
	"sharaga/pkg/logx"
	"time"
)

type Cache interface {
	SaveRefreshToken(ctx context.Context, refreshToken string, userId int, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, refreshToken string) (int, error)
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type Repository interface {
	GetWithPassword(ctx context.Context, username, password string) (*entity.User, error)
}

type Service struct {
	repo  Repository
	cache Cache

	cfg config.AuthConfig
}

func NewService(cfg config.AuthConfig, repo Repository, cache Cache) *Service {
	return &Service{
		cfg:   cfg,
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*entity.Tokens, error) {
	const op = "auth.Login"

	user, err := s.repo.GetWithPassword(ctx, email, password)
	if err != nil {
		if failure.IsNotFoundError(err) {
			return nil, failure.NewUnauthorizedError("invalid username or password")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, accessExpired, err := s.NewAccessToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken := s.newRefreshToken()
	refreshExpired := time.Now().Add(time.Duration(s.cfg.RefreshTokenTTL) * time.Second)
	if err := s.cache.SaveRefreshToken(ctx, refreshToken, user.Id, time.Duration(s.cfg.RefreshTokenTTL)*time.Second); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entity.Tokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessExpired,
		RefreshTokenExpiresAt: refreshExpired,
	}, nil
}

func (s *Service) UpdateTokens(ctx context.Context, refresh string) (*entity.Tokens, error) {
	const op = "auth.UpdateTokens"

	if len(refresh) == 0 {
		return nil, failure.NewInvalidRequestError("missig refresh token")
	}
	userId, err := s.cache.GetRefreshToken(ctx, refresh)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if userId == 0 {
		return nil, failure.NewInvalidRequestError("invalid refresh token")
	}

	if err := s.cache.DeleteRefreshToken(ctx, refresh); err != nil {
		contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "delete refresh token error", logx.Error(err), slog.Int("userId", userId))
	}

	access, accessExpired, err := s.NewAccessToken(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpired := time.Now().Add(time.Duration(s.cfg.RefreshTokenTTL) * time.Second)
	newRefresh := s.newRefreshToken()
	if err := s.cache.SaveRefreshToken(ctx, newRefresh, userId, time.Duration(s.cfg.RefreshTokenTTL)*time.Second); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entity.Tokens{
		AccessToken:           access,
		RefreshToken:          newRefresh,
		AccessTokenExpiresAt:  accessExpired,
		RefreshTokenExpiresAt: refreshExpired,
	}, nil
}

func (s *Service) DecodeAccessToken(access string) (int, error) {
	claims, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.AccessKey), nil
	})

	if err != nil {
		return 0, failure.NewUnauthorizedError(err.Error())
	}
	mapClaims, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return 0, failure.NewUnauthorizedError("claims is not a map")
	}

	id, ok := mapClaims["id"]
	if !ok {
		return 0, failure.NewUnauthorizedError("id is not in map claims")
	}
	tUnix, ok := mapClaims["expires"]
	if !ok {
		return 0, failure.NewUnauthorizedError("expires time is not in map claims")
	}

	if time.Now().After(time.Unix(int64(tUnix.(float64)), 0)) {
		return 0, failure.NewUnauthorizedError("token expired")
	}

	return int(id.(float64)), nil
}

func (s *Service) NewAccessToken(id int) (string, time.Time, error) {
	expires := time.Now().Add(time.Duration(s.cfg.AccessTokenTTL) * time.Second)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":      id,
		"expires": expires.Unix(),
	})
	token, err := claims.SignedString([]byte(s.cfg.AccessKey))
	if err != nil {
		return "", expires, failure.NewInternalError("create token error: " + err.Error())
	}
	return token, expires, nil
}

func (s *Service) newRefreshToken() string {
	token := uuid.NewString()
	return token
}
