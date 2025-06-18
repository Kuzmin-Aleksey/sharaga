package failure

import (
	"errors"
)

type UnauthorizedError struct {
}

func NewUnauthorizedError() UnauthorizedError {
	return UnauthorizedError{}
}

func (err UnauthorizedError) Error() string {
	return "Unauthorized"
}

func IsUnauthorizedError(err error) bool {
	return errors.As(err, new(UnauthorizedError))
}
