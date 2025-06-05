package server

import (
	"context"
)

type DB interface {
	RegisterUser(ctx context.Context, user *User) error
	CheckUser(ctx context.Context, user *User) error
	SaveMessage(ctx context.Context, message *Message) error
	GetLastMessage(ctx context.Context) (*Message, error)
	GetMessages(ctx context.Context, start uint64, end uint64) ([]Message, error)
}

type Message struct {
	MessageId uint64
	Username  string
	Content   string
}

type User struct {
	Username string
	Password string
}
