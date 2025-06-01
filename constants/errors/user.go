package errorConstants

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrPasswordIncorect    = errors.New("password incorect")
	ErrUsernameExist       = errors.New("username already exist")
	ErrEmailExist          = errors.New("email already exist")
	ErrPasswordDoesntMatch = errors.New("password doesn't match")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordIncorect,
	ErrUsernameExist,
	ErrPasswordDoesntMatch,
}
