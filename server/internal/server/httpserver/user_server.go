package httpserver

type userService interface {
}

type UserServer struct {
	userService userService
}

func NewUserServer(userService userService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}
