package main

import (
	"context"
	"github.com/soichisumi/go-util/logger"
	"github.com/soichisumi/grpc-echo-server/pkg/proto"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewEchoServer(endpoint string) *EchoServer {
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithUnaryInterceptor(),
	)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	c := grpctesting.NewEchoServiceClient(conn)

	return &EchoServer{
		client: c,
	}
}

type EchoServer struct{
	client grpctesting.EchoServiceClient
}


func (e *EchoServer) Echo(ctx context.Context, req *grpctesting.EchoRequest) (*grpctesting.EchoResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	logger.Info("", zap.Any("req", req), zap.Any("md", md))
	return e.client.Echo(ctx, req)
}

func (e *EchoServer) Empty(ctx context.Context, req *grpctesting.EmptyRequest) (*grpctesting.EmptyResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	logger.Info("", zap.Any("req", req), zap.Any("md", md))
	return e.client.Empty(ctx, req)
}
