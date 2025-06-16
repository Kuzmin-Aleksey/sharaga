package httpserver

import (
	"context"
	"net/http"
	"sharaga/pkg/contextx"
	"sharaga/pkg/failure"
	"strings"
)

type roleProvider interface {
	GetRole(ctx context.Context, userId int) (string, error)
}

type tokenDecoder interface {
	DecodeAccessToken(access string) (int, error)
}

type MwAuth struct {
	auth  tokenDecoder
	roles roleProvider
}

func NewMwAuth(auth tokenDecoder, roles roleProvider) *MwAuth {
	return &MwAuth{
		auth:  auth,
		roles: roles,
	}
}

func (a *MwAuth) withAuth(next http.HandlerFunc, roles ...string) http.HandlerFunc {
	rolesMap := make(map[string]bool)
	for _, role := range roles {
		rolesMap[role] = true
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		typeAndToken := strings.Split(r.Header.Get("Authorization"), " ")
		if len(typeAndToken) != 2 {
			writeAndLogErr(ctx, w, failure.NewUnauthorizedError("invalid token format"))
			return
		}
		if typeAndToken[0] != "Bearer" {
			writeAndLogErr(ctx, w, failure.NewUnauthorizedError("invalid auth type: "+typeAndToken[0]))
			return
		}
		token := typeAndToken[1]
		userId, err := a.auth.DecodeAccessToken(token)
		if err != nil {
			writeAndLogErr(ctx, w, err)
			return
		}

		role, err := a.roles.GetRole(ctx, userId)
		if err != nil {
			if failure.IsNotFoundError(err) {
				writeAndLogErr(ctx, w, failure.NewUnauthorizedError("role not found"))
			}

			writeAndLogErr(ctx, w, err)
			return
		}

		if !rolesMap[role] {
			writeAndLogErr(ctx, w, failure.NewUnauthorizedError("forbidden: "+role))
			return
		}

		r = r.WithContext(contextx.WithUserId(ctx, contextx.UserId(userId)))

		next.ServeHTTP(w, r)
	}
}
