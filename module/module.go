package module

import (
	"context"

	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/go-meowalien-lib/errs"

	"github.com/meowalien/RabbitGather-interest-crawler.git/lib"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler/factory"
)

// Module is the interface that contains all the methods that a module should have.
// any request from any connection will be handled by this interface.
type Module interface {
	interest.InterestCrawlerServer
	// Crawl creates a crawler and start crawling.
	Crawl(ctx context.Context, request *interest.CrawlRequest) (resp *interest.CrawlResponse, err error)
}

type Constructor struct {
	Factory factory.Factory
}

func (m Constructor) New() Module {
	return &module{
		factory: m.Factory,
	}
}

type module struct {
	interest.UnimplementedInterestCrawlerServer
	factory factory.Factory
}

func (s *module) Crawl(ctx context.Context, request *interest.CrawlRequest) (resp *interest.CrawlResponse, err error) {
	crawler, err := s.factory.NewCrawler(request.GetType(), request.GetMessage())
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
