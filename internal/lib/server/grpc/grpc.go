package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/meowalien/RabbitGather-golang.git/internal/lib/config"
	"github.com/meowalien/RabbitGather-golang.git/internal/lib/errs"
	"github.com/meowalien/RabbitGather-golang.git/internal/lib/graceful_shutdown"
)

func New(grpcConf config.GRPCConfig) *grpc.Server {
	serverParameters := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
	enforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	return grpc.NewServer(grpc.KeepaliveEnforcementPolicy(enforcementPolicy), grpc.KeepaliveParams(serverParameters))
}

func ListenAndServe(gracefulShutdown graceful_shutdown.GracefulShutdown, server *grpc.Server, port uint32) (err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		err = errs.New(err)
		return
	}
	gracefulShutdown.Add("grpcServer-listener", func(ctx context.Context) {
		err = listener.Close()
		if err != nil {
			slog.Error(errs.New(err).Error())
		}
	})

	gracefulShutdown.Add("grpcServer", func(ctx context.Context) {
		server.GracefulStop()
	})
	err = server.Serve(listener)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
