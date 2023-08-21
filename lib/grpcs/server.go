package grpcs

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/meowalien/go-meowalien-lib/errs"
)

// GRPCServerConstructor is a constructor for grpc.Server
// If the ServerParameters or EnforcementPolicy is nil, it will be set to default value
type GRPCServerConstructor struct {
	*keepalive.ServerParameters
	*keepalive.EnforcementPolicy
}

func (g GRPCServerConstructor) New() (grpcServer *grpc.Server) {
	if g.ServerParameters == nil {
		g.ServerParameters = &keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
			MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
			MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
			Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
		}
	}
	if g.EnforcementPolicy == nil {
		g.EnforcementPolicy = &keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}
	}

	grpcServer = grpc.NewServer(grpc.KeepaliveEnforcementPolicy(*g.EnforcementPolicy), grpc.KeepaliveParams(*g.ServerParameters))
	return grpcServer
}

// GRPCListenAndServeLauncher is a launcher for grpc.Server
// If the Port is 0, it will be set to default value 50051
type GRPCListenAndServeLauncher struct {
	GRPCServer *grpc.Server
	Port       uint32
}

func (g *GRPCListenAndServeLauncher) Name() string {
	return "GRPCServer"
}

func (g *GRPCListenAndServeLauncher) GracefulStop(ctx context.Context) (err error) {
	g.GRPCServer.GracefulStop()
	return
}

func (g *GRPCListenAndServeLauncher) ListenAndServe() (err error) {
	if g.Port == 0 {
		g.Port = 50051
	}
	port := fmt.Sprintf(":%d", g.Port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		err = errs.New(err)
		return
	}
	err = g.GRPCServer.Serve(listener)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
