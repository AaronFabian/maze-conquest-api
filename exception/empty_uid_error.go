package exception

type EmptyUidError struct {
	Code    int
	Status  string
	Message string
}

func NewEmptyUidError() *EmptyUidError {
	return &EmptyUidError{
		Code:    400,
		Status:  "Bad Request",
		Message: "UID is not provided !",
	}
}

func (notFoundError *EmptyUidError) Error() string {
	return notFoundError.Message
}
