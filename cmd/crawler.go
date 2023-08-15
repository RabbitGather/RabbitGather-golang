// the cmd package is used to handle the interface between the business logic and dependencies.
package cmd

import (
	"github.com/meowalien/RabbitGather-interest-crawler.git/module/crawler"
	"github.com/meowalien/RabbitGather-interest-crawler.git/server"
)

func CrawlerServer() server.Service {
	return crawler.NewServer()
}
