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
		srv.Logger.Info("Could not encrypt", "password", password)
		return nil, status.Error(codes.Internal, "couldn't encrypt your password")
	}
	user := &User{logReq.GetUsername(), password}
	err = srv.Db.RegisterUser(ctx, user)
	if err != nil {
		if err == custom_errors.ErrUserAlreadyExist {
			srv.Logger.Info("Trying to register already existing", "username", user.Username)
			return nil, status.Error(codes.AlreadyExists, "user with this name already exists")
		} else {
			srv.Logger.Warn("Error in db.RegisterUser", "error", err.Error())
			return nil, status.Error(codes.Internal, "could't save your account")
		}
	}

	srv.Logger.Info("Registered", "user", user.Username)
	return &proto.LoginResponse{}, nil
}

func (srv *GrpcServer) Stream(cnn proto.Chat_StreamServer) error {
	ctx := cnn.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		srv.Logger.Warn("could not load metadata")
		return status.Error(codes.Unavailable, "can not load metadata")
	}
	usernames := md.Get("username")
	if len(usernames) != 1 {
		srv.Logger.Warn("invalid request to Stream", "len", len(usernames))
		return status.Error(codes.Unauthenticated, "Invalid username in metadata")
	}
	username := usernames[0]
	srv.Connections.RegisterNewUser(username, cnn)

	message, err := srv.Db.GetLastMessage(ctx)
	if err == custom_errors.ErrNoMessages {
		srv.Logger.Info("no messages currently, send to username")
		message = &Message{}
	} else if err != nil {
		srv.Logger.Warn("could not get last message", "err", err.Error())
		return status.Error(codes.Internal, "could not get last message")
	}

	err = cnn.Send(&proto.Message{MessageId: message.MessageId, Username: message.Username, Content: message.Content})
	if err != nil {
		srv.Logger.Info("could not send last message to", "username", username)
		return status.Error(codes.Unavailable, "can not send last message")
	}

	for {
		sendReq, err := cnn.Recv()
		if err != nil {
			srv.Logger.Info("close connection", "username", username)
			return status.Error(codes.Aborted, "lost connection")
		}
		id, err := srv.Db.SaveMessage(ctx, &Message{Username: username, Content: sendReq.GetText()})
		if err != nil {
			srv.Logger.Warn("could not save message", "err", err.Error())
			continue
		}
		srv.Connections.SendMessageToActiveUsers(&proto.Message{MessageId: id, Username: username, Content: sendReq.GetText()})
	}
}

func (srv *GrpcServer) LoadHistory(ctx context.Context, hisReq *proto.HistoryRequest) (*proto.HistoryResponse, error) {
	start := int64(hisReq.LastMessageId) - int64(hisReq.Amount)
	start = max(start, 1)
	if hisReq.LastMessageId < 1 {
		srv.Logger.Warn("Invalid request at LoadHistory", "LastMessageId", hisReq.LastMessageId)
		return nil, status.Error(codes.InvalidArgument, "last_message id should be >= 1")
	}

	messages, err := srv.Db.GetMessages(ctx, uint64(start), hisReq.LastMessageId)
	if err != nil {
		srv.Logger.Warn("Error in db.GetMessages", "start", start, "end", hisReq.LastMessageId)
		return nil, status.Error(codes.Internal, "could not load history of messages")
	}
	var protoMessages []*proto.Message
	for _, v := range messages {
		protoMessages = append(protoMessages, &proto.Message{MessageId: v.MessageId, Username: v.Username, Content: v.Content})
	}
	srv.Logger.Info("Succesfully loaded history")
	return &proto.HistoryResponse{Messages: protoMessages}, nil
}
