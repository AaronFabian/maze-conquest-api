package exception

type EmptyUidError struct {
	Code    int
	Message string
}

func NewEmptyUidError(code int, error string) *EmptyUidError {
	return &EmptyUidError{
		Code:    code,
		Message: error,
	}
}

func (notFoundError *EmptyUidError) Error() string {
	return notFoundError.Message
}
