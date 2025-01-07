package exception

type NotFoundError struct {
	Code    int
	Message string
}

func NewNotFoundError(code int, error string) *NotFoundError {
	return &NotFoundError{
		Code:    code,
		Message: error,
	}
}

func (notFoundError *NotFoundError) Error() string {
	return notFoundError.Message
}
