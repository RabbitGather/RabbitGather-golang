package server

import "google.golang.org/grpc"

type GRPCService interface {
	MountRPCServer(grpcServer *grpc.Server)
}
