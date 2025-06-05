package storage

import (
	"context"
	"database/sql"

	"github.com/teryble09/go_grpc_chat/server"
	"github.com/teryble09/go_grpc_chat/server/custom_errors"
)

var stmtSaveMessage *sql.Stmt

func (db dbPostgre) SaveMessage(ctx context.Context, message *server.Message) (uint64, error) {
	res, err := stmtSaveMessage.ExecContext(ctx, message.Username, message.Content)
	if err != nil {
		return 0, custom_errors.ErrMessageFailedSave
	}
	id, _ := res.LastInsertId()
	return uint64(id), nil
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
	rows, err := stmtGetMessages.QueryContext(ctx, start, end)
	defer rows.Close()
	if err != nil {
		return nil, custom_errors.ErrMessageFailedRetrieve
	}

	for rows.Next() {
		var (
			id      uint64
			sender  string
			content string
		)
		err := rows.Scan(&id, &sender, &content)
		if err != nil {
			return nil, custom_errors.ErrMessageFailedRetrieve
		}
		messages = append(messages, server.Message{MessageId: id, Username: sender, Content: content})
	}
	if rows.Err() != nil {
		return nil, custom_errors.ErrMessageFailedRetrieve
	}

	return messages, nil
}

var stmtGetLastMessage *sql.Stmt

func (db dbPostgre) GetLastMessage(ctx context.Context) (*server.Message, error) {
	var message server.Message
	row := stmtGetLastMessage.QueryRowContext(ctx)
	err := row.Scan(&message.MessageId, &message.Username, &message.Content)
	if err != nil {
		return nil, custom_errors.ErrMessageFailedRetrieve
	}
	return &message, nil
}
