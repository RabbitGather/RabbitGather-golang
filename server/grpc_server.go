package server

import (
	"context"
	"fmt"

	"github.com/meowalien/go-meowalien-lib/grpcs"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	Port       int
	servers    *Servers
	grpcServer *grpc.Server
}

func (s *GRPCServer) AddService(grpcService GRPCService) {
	grpcService.MountTo(s.grpcServer)
}

func (s *GRPCServer) GracefulStop(ctx context.Context) (err error) {
	s.grpcServer.GracefulStop()
	return
}

func (s *GRPCServer) ListenAndServe() (err error) {
	if s.Port == 0 {
		return
	}
	err = grpcs.ListenAndServe(s.grpcServer, fmt.Sprintf(":%d", s.Port))
	if err != nil {
		fmt.Println("grpcServer stopped with error: ", err)
		return err
	}
	return
}
