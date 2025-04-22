package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/teryble09/go_grpc_chat/proto"
	"github.com/teryble09/go_grpc_chat/server"
	"google.golang.org/grpc"
)

func main() {
	log := slog.Default()

	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Error("Failed to listen: " + err.Error())
		os.Exit(1)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterChatServer(grpcServer, &server.GrpcServer{})

	log.Info("Starting server on port " + strconv.Itoa(port))
	grpcServer.Serve(lis)
}
