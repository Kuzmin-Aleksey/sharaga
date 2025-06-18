package failure

import (
	"errors"
	"fmt"
)

type ServerError struct {
	Msg  string
	Code int
}

func NewServerError(msg string, code int) ServerError {
	return ServerError{
		Msg:  msg,
		Code: code,
	}
}

func (err ServerError) Error() string {
	return fmt.Sprintf("Error: %d %s", err.Code, err.Msg)
}

func IsServerError(err error) bool {
	return errors.As(err, new(ServerError))
}
