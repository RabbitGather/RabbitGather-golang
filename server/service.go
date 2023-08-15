package server

import "google.golang.org/grpc"

type Service struct {
	GRPCService GRPCService
}

type GRPCService interface {
	MountTo(grpcServer *grpc.Server)
}
