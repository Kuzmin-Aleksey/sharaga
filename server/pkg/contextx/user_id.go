package contextx

import (
	"context"
	"strconv"
)

type UserId int

func (id UserId) String() string {
	return strconv.Itoa(int(id))
}

type contextKeyUserId struct{}

func WithUserId(ctx context.Context, userId UserId) context.Context {
	return context.WithValue(ctx, contextKeyUserId{}, userId)
}

func GetUserId(ctx context.Context) UserId {
	v, _ := ctx.Value(contextKeyUserId{}).(UserId)
	return v
}
