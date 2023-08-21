package main

import (
	"github.com/meowalien/go-meowalien-lib/graceful_shutdown"
	"github.com/meowalien/go-meowalien-lib/schedule"

	"github.com/meowalien/RabbitGather-interest-crawler.git/fremwork/connect"
	"github.com/meowalien/RabbitGather-interest-crawler.git/fremwork/server"
	"github.com/meowalien/RabbitGather-interest-crawler.git/fremwork/server/manager"
	"github.com/meowalien/RabbitGather-interest-crawler.git/lib/grpcs"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module/adepter"
	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler/factory"
)

func main() {
	grpcServer := grpcs.GRPCServerConstructor{}.New()
	cp := connect.ChannelPoolConstructor{
		Connection: connect.ConnectionConstructor{
			UserName: "rabbit",
			Password: "aneMicC7A9np",
			Address:  "localhost:5672",
		},
		MaxSize: 10,
	}.New()
	{
		md := module.Constructor{
			Factory: factory.Constructor{}.New(),
		}.New()
		cs := adepter.Constructor{
			Module: md,
		}.New()
		cs.GRPC(grpcServer)
		cs.RabbitMQ(cp)
	}
	gracefulShutdown := graceful_shutdown.NewGracefulShutdown()
	err := manager.Constructor{
		Retryer:          schedule.Retryer{},
		Servers:          []server.Launcher{&grpcs.GRPCListenAndServeLauncher{Port: 50051, GRPCServer: grpcServer}},
		GracefulShutdown: gracefulShutdown,
	}.New().ListenAndServe()
	if err != nil {
		panic(err)
	}
}
