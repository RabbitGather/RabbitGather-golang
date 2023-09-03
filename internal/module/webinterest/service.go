package webinterest

import (
	"context"

	"github.com/meowalien/RabbitGather-golang.git/proto/pb/webinterest"
)

type Service interface {
}

type service struct {
	webinterest.UnimplementedInterestCrawlerServer
}

func (g *service) Crawl(ctx context.Context, request *webinterest.CrawlRequest) (*webinterest.CrawlResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g *service) mustEmbedUnimplementedInterestCrawlerServer() {
	//TODO implement me
	panic("implement me")
}
