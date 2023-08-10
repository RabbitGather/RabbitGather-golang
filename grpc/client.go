package grpc

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func StartGRPCClient() *grpc.ClientConn {
	// TLS configuration
	tlsConfig := &tls.Config{
		// If you have server's certificate for verification:
		// RootCAs: yourCertPool,
		InsecureSkipVerify: true, // NOTE: Use this for testing purposes only. Always verify server's certificate in production.
	}

	// Keepalive parameters
	kaOpts := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial the gRPC server with timeout, TLS, and keepalive
	conn, err := grpc.DialContext(
		ctx,
		grpcPort,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithKeepaliveParams(kaOpts),
		grpc.WithBlock(), // block until the connection is established
	)
	if err != nil {
		log.Fatalf("Could not connect to %s: %v", grpcPort, err)
	}
	//defer conn.Close()
	return conn
	// Now you can use `conn` to create client stubs and make RPC calls.
	// ...
}
