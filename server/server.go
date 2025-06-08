package server

import (
	"context"
	"log/slog"

	"github.com/teryble09/go_grpc_chat/proto"
	"github.com/teryble09/go_grpc_chat/server/custom_errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	proto.UnimplementedChatServer
	Connections ConnStorage
	Db          DB
	Logger      *slog.Logger
}

func (srv *GrpcServer) Login(ctx context.Context, logReq *proto.LoginRequest) (*proto.LoginResponse, error) {
	password, err := HashPassword(logReq.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't encrypt your password")
	}
	user := &User{logReq.GetUsername(), password}
	err = srv.Db.RegisterUser(ctx, user)
	if err != nil {
		if err == custom_errors.ErrUserAlreadyExist {
			return nil, status.Error(codes.AlreadyExists, "user with this name already exists")
		} else {
			return nil, status.Error(codes.Internal, "could't save your account")
		}
	}

	return &proto.LoginResponse{}, nil
}

func (srv *GrpcServer) Stream(cnn proto.Chat_StreamServer) error {
	ctx := cnn.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unavailable, "can not load metadata")
	}
	usernames := md.Get("username")
	if len(usernames) != 1 {
		return status.Error(codes.Unauthenticated, "Invalid username in metadata")
	}
	username := usernames[0]
	srv.Connections.RegisterNewUser(username, cnn)

	message, err := srv.Db.GetLastMessage(ctx)
	if err != nil {
		return status.Error(codes.Internal, "could not get last message")
	}

	err = cnn.Send(&proto.Message{MessageId: message.MessageId, Username: message.Username, Content: message.Content})
	if err != nil {
		return status.Error(codes.Unavailable, "can not send last message")
	}

	for {
		sendReq, err := cnn.Recv()
		if err != nil {
			return status.Error(codes.Aborted, "lost connection")
		}
		id, err := srv.Db.SaveMessage(ctx, &Message{Username: username, Content: sendReq.GetText()})
		if err != nil {
			//log
			continue
		}
		srv.Connections.SendMessageToActiveUsers(&proto.Message{MessageId: id, Username: username, Content: sendReq.GetText()})
	}
}

func (srv *GrpcServer) LoadHistory(ctx context.Context, hisReq *proto.HistoryRequest) (*proto.HistoryResponse, error) {
	start := int64(hisReq.LastMessageId) - int64(hisReq.Amount)
	if start < 1 {
		start = 1
	}
	if hisReq.LastMessageId < 1 {
		return nil, status.Error(codes.InvalidArgument, "last_message id should be >= 1")
	}

	messages, err := srv.Db.GetMessages(ctx, uint64(start), hisReq.LastMessageId)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not load history of messages")
	}
	var protoMessages []*proto.Message
	for _, v := range messages {
		protoMessages = append(protoMessages, &proto.Message{MessageId: v.MessageId, Username: v.Username, Content: v.Content})
	}
	return &proto.HistoryResponse{Messages: protoMessages}, nil
}
