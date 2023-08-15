package server

import (
	"context"

	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/grpcs"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	GRPCServer                 *grpc.Server
	GRPCListenAndServeLauncher grpcs.GRPCListenAndServeLauncher
	allServiceMount            []GRPCService
}

func (s *GRPCServer) Name() string {
	return "GRPCServer"
}

func (s *GRPCServer) GracefulStop(ctx context.Context) (err error) {
	s.GRPCServer.GracefulStop()
	return
}

func (s *GRPCServer) ListenAndServe() (err error) {
	for _, grpcService := range s.allServiceMount {
		grpcService.MountRPCServer(s.GRPCServer)
	}

	err = s.GRPCListenAndServeLauncher.ListenAndServe(s.GRPCServer)
	if err != nil {
		err = errs.New(err)
		return err
	}
	return
}

// AddService store GRPCService , and will call MountRPCServer() when ListenAndServe()
func (s *GRPCServer) AddService(cs GRPCService) {
	s.allServiceMount = append(s.allServiceMount, cs)
}
