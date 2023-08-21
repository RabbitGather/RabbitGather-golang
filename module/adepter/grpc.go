package adepter

import (
	"github.com/meowalien/RabbitGather-proto/go/interest"
	"google.golang.org/grpc"
)

func (a *adepter) GRPC(svc *grpc.Server) {
	interest.RegisterInterestCrawlerServer(svc, a.md)
}
