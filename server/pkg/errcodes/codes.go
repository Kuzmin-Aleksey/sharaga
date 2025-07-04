package errcodes

type Code string

func (e Code) String() string {
	return string(e)
}

const (
	Internal       Code = "InternalError"
	NotFound       Code = "NotFound"
	InvalidRequest Code = "InvalidRequest"
	Unauthorized   Code = "Unauthorized"
)
