package custom_errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserWrongPassword = errors.New("wrong password")
	ErrUserAlreadyExist  = errors.New("user already exists")
)
