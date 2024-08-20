package main

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-nashville/internal/config"
	"github.com/bhupeshpandey/task-manager-nashville/internal/logger"
	"github.com/bhupeshpandey/task-manager-nashville/internal/metrics"
	grpcpkg "github.com/bhupeshpandey/task-manager-nashville/internal/proto"
	"github.com/bhupeshpandey/task-manager-nashville/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// Load configuration
	conf, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create logger based on configuration
	newLogger := logger.NewLogger(conf.LoggingConfig)

	// Update metrics config to use logger and create the prometheus registry.
	conf.MetricsConfig.Logger = newLogger
	metrics.InitializePrometheus(conf.MetricsConfig)

	// Create GRPC client to connect to the gallatin service grpc server.
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", conf.GrpcServer.Host, conf.GrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcClient := grpcpkg.NewTaskServiceClient(conn)

	conf.HttpServer.Logger = newLogger
	httpServer := server.NewHTTPServer(grpcClient, conf.HttpServer)

	server.Start(httpServer)
}
