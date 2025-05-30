package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/teryble09/go_grpc_chat/server"
	"github.com/teryble09/go_grpc_chat/server/custom_errors"
)

var stmtRegisterUser *sql.Stmt

func (db dbPostgre) RegisterUser(ctx context.Context, user *server.User) error {
	_, err := stmtRegisterUser.ExecContext(ctx, user.Username, user.Password)
	if err != nil {
		return custom_errors.ErrUserAlreadyExist
	}

	return nil
}

var stmtGetPassword *sql.Stmt

func (db dbPostgre) CheckUser(ctx context.Context, user *server.User) error {
	res := stmtGetPassword.QueryRowContext(ctx, user.Username)

	var password string
	err := res.Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return custom_errors.ErrUserNotFound
		}
		return err
	}

	if password != user.Password {
		return custom_errors.ErrUserWrongPassword
	}

	return nil
}
