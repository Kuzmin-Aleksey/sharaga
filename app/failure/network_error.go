package failure

import (
	"errors"
)

type NetworkError struct {
	err error
}

func NewNetworkError(err error) NetworkError {
	return NetworkError{err: err}
}

func (e NetworkError) Error() string {
	return "NetworkError: " + e.err.Error()
}

func IsNetworkError(err error) bool {
	return errors.As(err, new(NetworkError))
}
