package webinterest

import (
	"google.golang.org/grpc"

	"github.com/meowalien/RabbitGather-golang.git/proto/pb/webinterest"
)

func New(grpcServer *grpc.Server) Service {
	svc := service{}
	webinterest.RegisterInterestCrawlerServer(grpcServer, &svc)
	return svc
}
