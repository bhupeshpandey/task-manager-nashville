package main

import (
	"fmt"
	"github.com/bhupeshpandey/task-manager-nashville/internal/config"
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

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", conf.GrpcServer.Host, conf.GrpcServer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcClient := grpcpkg.NewTaskServiceClient(conn)

	httpServer := server.NewHTTPServer(grpcClient)

	server.Start(httpServer)
}
