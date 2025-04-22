package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/teryble09/go_grpc_chat/server"
	"github.com/teryble09/go_grpc_chat/server/custom_errors"
)

func (db dbPostgre) RegisterUser(ctx context.Context, user *server.User) error {
	stmt, err := db.conn.PrepareContext(ctx, `
		INSERT INTO users (Username, HashPassword)
		VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, user.Username, user.Password)
	if err != nil {
		return custom_errors.ErrUserAlreadyExist
	}

	return nil
}

func (db dbPostgre) CheckUser(ctx context.Context, user *server.User) error {
	stmt, err := db.conn.PrepareContext(ctx, `
		SELECT HashPassword FROM users
		WHERE Username = $1
	`)
	if err != nil {
		return err
	}

	res := stmt.QueryRowContext(ctx, user.Username)

	var password string
	err = res.Scan(&password)
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
