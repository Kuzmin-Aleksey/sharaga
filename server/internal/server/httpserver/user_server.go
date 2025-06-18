package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sharaga/internal/domain/entity"
	"sharaga/pkg/contextx"
	"sharaga/pkg/failure"
	"sharaga/pkg/rest"
	"strconv"
)

type userService interface {
	NewUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	GetAll(ctx context.Context) ([]entity.User, error)
	DeleteUser(ctx context.Context, userId int) error
	GetById(ctx context.Context, userId int) (*entity.User, error)
}

type UserServer struct {
	userService userService
}

func NewUserServer(userService userService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

func (s *UserServer) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := new(entity.User)

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.userService.NewUser(ctx, user); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, rest.IdResponse{
		Id: user.Id,
	}, http.StatusOK)
}

func (s *UserServer) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := new(entity.User)

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.userService.UpdateUser(ctx, user); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserServer) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError("invalid user id: "+r.FormValue("user_id")))
		return
	}

	if err := s.userService.DeleteUser(ctx, userId); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserServer) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := s.userService.GetAll(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, users, http.StatusOK)
}

func (s *UserServer) GetSelf(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := contextx.GetUserId(ctx)

	user, err := s.userService.GetById(ctx, int(userId))
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, user, http.StatusOK)
}
