package grpc_interest

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/meowalien/RabbitGather-proto/proto/interest"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

const grpcServerListen = ":50051"

func TestCrawlerServer(t *testing.T) {
	listener, err := net.Listen("tcp", grpcServerListen)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcServerListen, err)
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	grpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))

	interest.RegisterInterestCrawlerServer(grpcServer, &InterestCrawlerServerImpl{})

	fmt.Println("Start listening on port", grpcServerListen)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

const grpcClientDial = "localhost:50051"

func TestCrawlerClient(t *testing.T) {

	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		grpcClientDial,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithBlock(), // block until the connection is established
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := interest.NewInterestCrawlerClient(conn)
	crawl, err := client.Crawl(context.Background(), &interest.CrawlRequest{
		Url:           "https://example.com",
		QuerySelector: ".some-class",
	})
	if err != nil {
		return
	}
	fmt.Println(crawl)
}
