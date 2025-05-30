package storage

import (
	"context"
	"database/sql"

	"github.com/teryble09/go_grpc_chat/server"
)

var stmtSaveMessage *sql.Stmt

func (db dbPostgre) SaveMessage(ctx context.Context, message *server.Message) error {
	_, err := stmtSaveMessage.ExecContext(ctx, message.Username, message.Content)
	return err
}

var stmtGetMessages *sql.Stmt

func (db dbPostgre) GetMessages(ctx context.Context, start uint64, end uint64) ([]server.Message, error) {
	var messages []server.Message
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	rows := stmtGetMessages.QueryRowContext(ctx, start, end)

	return messages, nil
}

var stmtGetLastMessage *sql.Stmt

func (db dbPostgre) GetLastMessage(ctx context.Context) (*server.Message, error) {
	var message server.Message
	rows := stmtGetLastMessage.QueryRowContext(ctx)

	return message, nil
}
