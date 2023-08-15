package crawler

import (
	"context"

	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/go-meowalien-lib/errs"
	"google.golang.org/grpc"

	"github.com/meowalien/RabbitGather-interest-crawler.git/lib"
)

type CrawlerService struct {
	interest.UnimplementedInterestCrawlerServer
	CrawlerFactory Factory
}

func (c *CrawlerService) MountRPCServer(grpcServer *grpc.Server) {
	interest.RegisterInterestCrawlerServer(grpcServer, c)
}

func (c *CrawlerService) Crawl(ctx context.Context, request *interest.CrawlRequest) (resp *interest.CrawlResponse, err error) {
	resp = &interest.CrawlResponse{}
	m := request.GetMessage()
	if m == nil {
		err = errs.New("message is nil")
		return nil, err
	}

	crawler, err := c.CrawlerFactory.NewCrawler(request.GetType(), m)
	if err != nil {
		err = errs.New(err)
		return nil, err
	}
	responseChan, errorChan := crawler.Crawl(ctx)

	err = lib.Pipe(responseChan, errorChan, func(input []byte) bool {
		if input == nil {
			return false
		}
		resp.Data = append(resp.Data, input)
		return true
	})
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
