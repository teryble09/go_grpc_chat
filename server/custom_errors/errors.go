package custom_errors

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserWrongPassword = errors.New("wrong password")
	ErrUserAlreadyExist  = errors.New("user already exists")

	ErrMessageFailedSave     = errors.New("can not save message")
	ErrMessageFailedRetrieve = errors.New("could not retrieve message")
	ErrNoMessages            = errors.New("no messages")
)
