package server

import (
	"context"

	"github.com/teryble09/go_grpc_chat/proto"
	"github.com/teryble09/go_grpc_chat/server/custom_errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	proto.UnimplementedChatServer
	Connections ConnStorage
	Db          DB
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

func (srv *GrpcServer) Stream(proto.Chat_StreamServer) error {
	return nil
}

func (srv *GrpcServer) LoadHistory(context.Context, *proto.HistoryRequest) (*proto.HistoryResponse, error) {
	return nil, nil
}
