package storage

import (
	"context"

	"github.com/teryble09/go_grpc_chat/server"
)

func (db dbPostgre) SaveMessage(ctx context.Context, message *server.Message) error {
	return nil
}

func (db dbPostgre) GetMessages(ctx context.Context, start uint64, end uint64) ([]server.Message, error) {
	return nil, nil
}
