package failure

type InvalidResponseError struct {
	err error
}

func NewInvalidResponseError(err error) NetworkError {
	return NetworkError{err: err}
}

func (e InvalidResponseError) Error() string {
	return "InvalidResponseError: " + e.err.Error()
}
