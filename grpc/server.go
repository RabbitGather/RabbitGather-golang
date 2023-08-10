package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

const grpcPort = "0.0.0.0:50051"

func StartGRPCServer(grpcServer *grpc.Server) {

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return
}
