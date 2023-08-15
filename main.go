package main

import (
	"github.com/meowalien/go-meowalien-lib/graceful_shutdown"

	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler"
	"github.com/meowalien/RabbitGather-interest-crawler.git/server"
)

func main() {
	gracefulShutdown := graceful_shutdown.NewGracefulShutdown()
	servers := server.Servers{
		GracefulShutdown: gracefulShutdown,
	}
	servers.AddService(crawler.NewServer())
}
