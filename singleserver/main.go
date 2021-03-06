package main

import (
	"flag"
	"fmt"
	"github.com/soichisumi-sandbox/opencensus-grpc-go-sample/opencensus"
	"github.com/soichisumi/go-util/logger"
	"github.com/soichisumi/grpc-echo-server/pkg/health"
	grpctesting "github.com/soichisumi/grpc-echo-server/pkg/proto"
	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const (
	defaultPort = 8080
)

func main() {
	var port = flag.Int("p", defaultPort, "port number for listening")
	flag.Parse()

	logger.Info("server is running.", zap.Int("port", *port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
	//creds, err := credentials.NewServerTLSFromFile("./certs/cert.pem", "./certs/privkey.pem")
	//if err != nil {
	//	logger.Fatal(err.Error(), zap.Error(err))
	//}
	//server := grpc.NewServer(grpc.Creds(creds))

	if os.Getenv("GCP_PROJECT") == "" {
		logger.Fatal("environment variable GCP_PROJECT is required")
	}

	ex := opencensus.SetupExporter(os.Getenv("GCP_PROJECT"))
	defer ex.Flush()

	server := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: true}),
		grpc.UnaryInterceptor(opencensus.UnaryServerTraceInterceptor()),
	)
	grpctesting.RegisterEchoServiceServer(server, NewEchoServer())
	grpc_health_v1.RegisterHealthServer(server, health.NewHealthServer())
	reflection.Register(server)

	logger.Info("", zap.Int("port", *port))
	if err := server.Serve(lis); err != nil {
		logger.Fatal(err.Error(), zap.Error(err))
	}
}


