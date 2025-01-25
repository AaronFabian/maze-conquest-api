package helper

import (
	"fmt"
	"maze-conquest-api/exception"
	"strings"

	"firebase.google.com/go/auth"
)

func PanicForGetAuth(err error, uid string) {
	if err != nil {
		switch {
		case auth.IsUserNotFound(err):
			panic(exception.NewNotFoundError("problem while updating, User with Uid " + uid + " not found"))
		case strings.Contains(err.Error(), "PERMISSION_DENIED"):
			err = fmt.Errorf("insufficient permissions to access user data: %w", err)
			panic(err)

		case strings.Contains(err.Error(), "FAILED_PRECONDITION"):
			err = fmt.Errorf("firebase authentication not properly initialized: %w", err)
			panic(err)

		default:
			err = fmt.Errorf("unexpected error getting user %s: %w", uid, err)
			panic(err)

		}
	}
}
