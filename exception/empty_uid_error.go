package exception

type EmptyUidError struct {
	Code    int
	Message string
}

func NewEmptyUidError() *EmptyUidError {
	return &EmptyUidError{
		Code:    400,
		Message: "UID is not provided !",
	}
}

func (notFoundError *EmptyUidError) Error() string {
	return notFoundError.Message
}
