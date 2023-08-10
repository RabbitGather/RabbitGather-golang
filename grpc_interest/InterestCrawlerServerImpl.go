package grpc_interest

import (
	"context"
	"fmt"

	"github.com/meowalien/RabbitGather-proto/proto/interest"
	"google.golang.org/protobuf/types/known/emptypb"
)

type InterestCrawlerServerImpl struct {
	interest.UnimplementedInterestCrawlerServer
}

func (s *InterestCrawlerServerImpl) Crawl(ctx context.Context, req *interest.CrawlRequest) (*emptypb.Empty, error) {
	fmt.Println("req: ", req)

	// Your implementation here.
	return &emptypb.Empty{}, nil
}
