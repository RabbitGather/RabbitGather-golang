package main

import (
	"github.com/meowalien/go-meowalien-lib/graceful_shutdown"
	"github.com/meowalien/go-meowalien-lib/grpcs"
	"github.com/meowalien/go-meowalien-lib/schedule"

	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler"
	"github.com/meowalien/RabbitGather-interest-crawler.git/server"
)

func main() {
	grpcServer := &server.GRPCServer{
		GRPCServer:                 grpcs.GRPCServerConstructor{}.New(),
		GRPCListenAndServeLauncher: grpcs.GRPCListenAndServeLauncher{},
	}

	cs := &crawler.CrawlerService{
		CrawlerFactory: crawler.CrawlerFactoryConstructor{}.New(),
	}
	grpcServer.AddService(cs)

	gracefulShutdown := graceful_shutdown.NewGracefulShutdown()
	err := server.ManagerConstructor{
		Retryer:          schedule.Retryer{},
		Servers:          []server.Server{grpcServer},
		GracefulShutdown: gracefulShutdown,
	}.New().ListenAndServe()
	if err != nil {
		panic(err)
	}
}
